import unittest

from zneondrive.domain import (
    BuildConflictError,
    CapacityError,
    DomainError,
    Entitlements,
    RaceValidationError,
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

    def test_quest_reward_is_idempotent(self) -> None:
        first = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        second = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        self.assertEqual(first, second)
        character = self.world.characters["char-1"]
        self.assertEqual((character.money, character.xp, character.reputation), (500, 100, 5))

    def test_same_quest_with_new_operation_does_not_double_reward(self) -> None:
        self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-1"
        )
        receipt = self.world.complete_quest(
            "acct-1", "MQ011", money=500, xp=100, reputation=5, operation_id="q-2"
        )
        self.assertEqual(receipt["money"], 0)
        self.assertEqual(self.world.characters["char-1"].money, 500)

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

    def test_starter_vehicle_cannot_be_routine_deleted(self) -> None:
        with self.assertRaises(DomainError):
            self.world.delete_vehicle("acct-1", "veh-1")


if __name__ == "__main__":
    unittest.main()
