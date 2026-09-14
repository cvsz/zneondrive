import unittest

from zneondrive.domain import (
    BlueprintRequiredError,
    BuildConflictError,
    CapacityError,
    DomainError,
    Entitlements,
    InsufficientInventoryError,
    OperationConflictError,
    QuestOutOfOrderError,
    RaceValidationError,
    STARTER_REBUILD_BLUEPRINT,
    WorldState,
)


STARTER_PARTS = [
    "part_chassis_starter_prototype",
    "part_engine_ice_street_i",
    "part_ecu_legacy_zero",
]


class WorldStateTests(unittest.TestCase):
    def setUp(self) -> None:
        self.world = WorldState()
        self.world.create_account("acct-1", "char-1")
        self.vehicle = self.world.grant_vehicle(
            "acct-1",
            "veh-1",
            STARTER_PARTS,
            starter_lineage=True,
            operation_id="grant-starter-1",
        )

    def complete_prerequisites(self, quest_id: str) -> None:
        target = int(quest_id[2:])
        character = self.world.characters["char-1"]
        for number in range(1, target):
            prerequisite = f"MQ{number:03d}"
            if prerequisite not in character.completed_quests:
                self.world.complete_quest(
                    "acct-1",
                    prerequisite,
                    operation_id=f"prereq-{prerequisite}",
                )

    def test_free_account_has_one_vehicle_capacity(self) -> None:
        with self.assertRaises(CapacityError):
            self.world.grant_vehicle(
                "acct-1", "veh-2", STARTER_PARTS, operation_id="grant-extra"
            )

    def test_vip_increases_capacity_not_competitive_signature(self) -> None:
        vip = WorldState()
        vip.create_account(
            "acct-vip", "char-vip", entitlements=Entitlements(garage_slots=5)
        )
        vehicle = vip.grant_vehicle(
            "acct-vip", "veh-vip", STARTER_PARTS, operation_id="grant-vip"
        )
        self.assertEqual(
            vip.competitive_build_signature(vehicle.vehicle_id),
            self.world.competitive_build_signature(self.vehicle.vehicle_id),
        )
        vip.grant_vehicle(
            "acct-vip", "veh-vip-2", STARTER_PARTS, operation_id="grant-vip-2"
        )
        self.assertEqual(len(vip.characters["char-vip"].vehicle_ids), 2)

    def test_build_revision_is_append_only_and_optimistic(self) -> None:
        old = self.vehicle.active_build
        new = self.world.revise_build(
            "acct-1",
            "veh-1",
            STARTER_PARTS + ["part_brakes_track_i"],
            expected_revision=1,
            operation_id="build-2",
        )
        self.assertEqual(old.revision, 1)
        self.assertEqual(new.revision, 2)
        self.assertEqual(self.vehicle.build_history[0], old)
        with self.assertRaises(BuildConflictError):
            self.world.revise_build(
                "acct-1",
                "veh-1",
                STARTER_PARTS,
                expected_revision=1,
                operation_id="stale-build",
            )

    def test_build_operation_replay_requires_same_semantic_payload(self) -> None:
        first = self.world.revise_build(
            "acct-1",
            "veh-1",
            STARTER_PARTS + ["part_brakes_track_i"],
            expected_revision=1,
            operation_id="build-replay",
        )
        replay = self.world.revise_build(
            "acct-1",
            "veh-1",
            list(reversed(STARTER_PARTS + ["part_brakes_track_i"])),
            expected_revision=1,
            operation_id="build-replay",
        )
        self.assertEqual(first, replay)
        self.assertEqual(self.vehicle.active_build.revision, 2)

        with self.assertRaises(OperationConflictError):
            self.world.revise_build(
                "acct-1",
                "veh-1",
                STARTER_PARTS + ["part_suspension_street_i"],
                expected_revision=1,
                operation_id="build-replay",
            )
        self.assertEqual(self.vehicle.active_build.revision, 2)

    def test_quest_ids_and_prerequisites_are_enforced(self) -> None:
        character = self.world.characters["char-1"]
        with self.assertRaises(DomainError):
            self.world.complete_quest(
                "acct-1", "MQ101", operation_id="invalid-quest"
            )
        with self.assertRaises(QuestOutOfOrderError):
            self.world.complete_quest(
                "acct-1", "MQ002", operation_id="skip-mq001"
            )
        self.assertEqual(character.completed_quests, set())

        self.world.complete_quest("acct-1", "MQ001", operation_id="mq001")
        self.world.complete_quest("acct-1", "MQ002", operation_id="mq002")
        self.assertEqual(character.completed_quests, {"MQ001", "MQ002"})

    def test_quest_reward_is_idempotent(self) -> None:
        self.complete_prerequisites("MQ011")
        first = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        second = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        self.assertEqual(first, second)
        character = self.world.characters["char-1"]
        self.assertEqual((character.money, character.xp, character.reputation), (500, 100, 5))

    def test_quest_operation_replay_rejects_changed_reward_payload(self) -> None:
        self.complete_prerequisites("MQ011")
        self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-conflict"
        )
        with self.assertRaises(OperationConflictError):
            self.world.complete_quest(
                "acct-1",
                "MQ011",
                money=5000,
                xp=100,
                reputation=5,
                operation_id="q-conflict",
            )
        character = self.world.characters["char-1"]
        self.assertEqual((character.money, character.xp, character.reputation), (500, 100, 5))

    def test_same_quest_with_new_operation_does_not_double_reward(self) -> None:
        self.complete_prerequisites("MQ011")
        self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        receipt = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-2"
        )
        self.assertEqual(receipt["money"], 0)
        self.assertEqual(self.world.characters["char-1"].money, 500)

    def test_garage_17_quest_side_effects_are_idempotent(self) -> None:
        character = self.world.characters["char-1"]
        self.complete_prerequisites("MQ004")
        self.world.complete_quest("acct-1", "MQ004", operation_id="mq004-a")
        self.assertEqual(character.inventory.get("part_brakes_track_i"), 1)
        self.world.complete_quest("acct-1", "MQ004", operation_id="mq004-b")
        self.assertEqual(character.inventory.get("part_brakes_track_i"), 1)

        self.world.complete_quest("acct-1", "MQ005", operation_id="mq005-a")
        self.assertIn(STARTER_REBUILD_BLUEPRINT, character.blueprints)
        self.world.complete_quest("acct-1", "MQ005", operation_id="mq005-b")
        self.assertEqual(character.blueprints, {STARTER_REBUILD_BLUEPRINT})

        self.complete_prerequisites("MQ009")
        self.world.complete_quest("acct-1", "MQ009", operation_id="mq009-a")
        self.world.complete_quest("acct-1", "MQ009", operation_id="mq009-b")
        self.assertEqual(character.inventory.get("part_tires_street_i"), 1)

    def test_rebuild_requires_blueprint_and_inventory_before_mutation(self) -> None:
        character = self.world.characters["char-1"]
        target = STARTER_PARTS + ["part_tires_street_i"]
        with self.assertRaises(BlueprintRequiredError):
            self.world.rebuild_vehicle(
                "acct-1",
                "veh-1",
                target,
                expected_revision=1,
                operation_id="rebuild-before-blueprint",
            )
        self.assertEqual(self.vehicle.active_build.revision, 1)
        self.assertEqual(character.inventory, {})

        self.complete_prerequisites("MQ005")
        self.world.complete_quest("acct-1", "MQ005", operation_id="mq005-unlock")
        inventory_before = dict(character.inventory)
        self.assertNotIn("part_tires_street_i", inventory_before)
        with self.assertRaises(InsufficientInventoryError):
            self.world.rebuild_vehicle(
                "acct-1",
                "veh-1",
                target,
                expected_revision=1,
                operation_id="rebuild-without-part",
            )
        self.assertEqual(self.vehicle.active_build.revision, 1)
        self.assertEqual(character.inventory, inventory_before)

    def test_rebuild_consumes_returns_and_replays_semantic_payload(self) -> None:
        character = self.world.characters["char-1"]
        self.complete_prerequisites("MQ004")
        self.world.complete_quest("acct-1", "MQ004", operation_id="mq004")
        self.world.complete_quest("acct-1", "MQ005", operation_id="mq005")
        first_parts = STARTER_PARTS + ["part_brakes_track_i"]
        first = self.world.rebuild_vehicle(
            "acct-1",
            "veh-1",
            first_parts,
            expected_revision=1,
            operation_id="rebuild-1",
        )
        self.assertEqual(first.revision, 2)
        self.assertNotIn("part_brakes_track_i", character.inventory)

        replay = self.world.rebuild_vehicle(
            "acct-1",
            "veh-1",
            list(reversed(first_parts)),
            expected_revision=1,
            operation_id="rebuild-1",
        )
        self.assertEqual(replay, first)
        self.assertEqual(self.vehicle.active_build.revision, 2)

        with self.assertRaises(OperationConflictError):
            self.world.rebuild_vehicle(
                "acct-1",
                "veh-1",
                STARTER_PARTS + ["part_tires_street_i"],
                expected_revision=1,
                operation_id="rebuild-1",
            )
        self.assertEqual(self.vehicle.active_build.revision, 2)

        self.complete_prerequisites("MQ009")
        self.world.complete_quest("acct-1", "MQ009", operation_id="mq009")
        swapped = self.world.rebuild_vehicle(
            "acct-1",
            "veh-1",
            STARTER_PARTS + ["part_tires_street_i"],
            expected_revision=2,
            operation_id="rebuild-2",
        )
        self.assertEqual(swapped.revision, 3)
        self.assertEqual(character.inventory.get("part_brakes_track_i"), 1)
        self.assertNotIn("part_tires_street_i", character.inventory)

    def test_vehicle_grant_operation_replay_rejects_changed_vehicle(self) -> None:
        vip = WorldState()
        vip.create_account(
            "acct-vip", "char-vip", entitlements=Entitlements(garage_slots=3)
        )
        first = vip.grant_vehicle(
            "acct-vip", "veh-a", STARTER_PARTS, operation_id="grant-replay"
        )
        replay = vip.grant_vehicle(
            "acct-vip", "veh-a", list(reversed(STARTER_PARTS)), operation_id="grant-replay"
        )
        self.assertIs(first, replay)
        with self.assertRaises(OperationConflictError):
            vip.grant_vehicle(
                "acct-vip", "veh-b", STARTER_PARTS, operation_id="grant-replay"
            )
        self.assertEqual(vip.characters["char-vip"].vehicle_ids, ["veh-a"])

    def test_race_result_is_bound_to_registered_build_and_checkpoint_order(self) -> None:
        registration = self.world.register_race(
            "acct-1", "race-1", "veh-1", ["start", "mid", "finish"]
        )
        result = self.world.submit_race_result(
            "acct-1",
            "race-1",
            build_revision=registration.accepted_build_revision,
            checkpoint_times_ms=[100, 5000, 10000],
            operation_id="race-result-1",
        )
        self.assertEqual(result.elapsed_ms, 10000)
        with self.assertRaises(RaceValidationError):
            self.world.submit_race_result(
                "acct-1",
                "race-1",
                build_revision=registration.accepted_build_revision,
                checkpoint_times_ms=[100, 5000, 4999],
                operation_id="race-result-bad",
            )

    def test_race_result_operation_replay_rejects_changed_timing_payload(self) -> None:
        registration = self.world.register_race(
            "acct-1", "race-replay", "veh-1", ["start", "mid", "finish"]
        )
        first = self.world.submit_race_result(
            "acct-1",
            "race-replay",
            build_revision=registration.accepted_build_revision,
            checkpoint_times_ms=[100, 5000, 10000],
            operation_id="race-replay-op",
        )
        replay = self.world.submit_race_result(
            "acct-1",
            "race-replay",
            build_revision=registration.accepted_build_revision,
            checkpoint_times_ms=[100, 5000, 10000],
            operation_id="race-replay-op",
        )
        self.assertEqual(first, replay)
        with self.assertRaises(OperationConflictError):
            self.world.submit_race_result(
                "acct-1",
                "race-replay",
                build_revision=registration.accepted_build_revision,
                checkpoint_times_ms=[100, 5000, 9999],
                operation_id="race-replay-op",
            )
        self.assertEqual(self.world.race_results["race-replay"], first)

    def test_operation_id_cannot_cross_mutation_types(self) -> None:
        self.complete_prerequisites("MQ011")
        self.world.complete_quest(
            "acct-1", "MQ011", money=1, operation_id="cross-kind-op"
        )
        with self.assertRaises(OperationConflictError):
            self.world.revise_build(
                "acct-1",
                "veh-1",
                STARTER_PARTS + ["part_brakes_track_i"],
                expected_revision=1,
                operation_id="cross-kind-op",
            )

    def test_starter_vehicle_cannot_be_routine_deleted(self) -> None:
        with self.assertRaises(DomainError):
            self.world.delete_vehicle("acct-1", "veh-1")


if __name__ == "__main__":
    unittest.main()
