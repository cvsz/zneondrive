import json
from pathlib import Path
import unittest

from zneondrive.domain import CapacityError, DomainError, Entitlements, WorldState


VECTORS_PATH = Path(__file__).with_name("entitlement-parity-vectors.json")
PARTS = ("part_chassis_starter_prototype", "part_engine_ice_street_i")


class EntitlementParityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.vectors = json.loads(VECTORS_PATH.read_text(encoding="utf-8"))
        if cls.vectors.get("schema_version") != 1:
            raise AssertionError("unsupported entitlement parity schema version")

    def test_fixture_integrity(self) -> None:
        garage_cases = self.vectors.get("garage_capacity_cases", [])
        delete_cases = self.vectors.get("routine_delete_cases", [])
        self.assertTrue(garage_cases, "garage_capacity_cases must not be empty")
        self.assertTrue(delete_cases, "routine_delete_cases must not be empty")
        self.assertEqual(len({case["name"] for case in garage_cases}), len(garage_cases))
        self.assertEqual(len({case["name"] for case in delete_cases}), len(delete_cases))
        for case in garage_cases:
            self.assertGreaterEqual(case["current_vehicles"], 0)
            self.assertGreaterEqual(case["garage_slots"], 1)

    def test_python_garage_capacity_matches_shared_vectors(self) -> None:
        for index, case in enumerate(self.vectors["garage_capacity_cases"], start=1):
            with self.subTest(vector=case["name"]):
                world = WorldState()
                account_id = f"acct-capacity-{index}"
                character_id = f"char-capacity-{index}"
                world.create_account(account_id, character_id, entitlements=Entitlements(garage_slots=case["garage_slots"]))
                for vehicle_index in range(case["current_vehicles"]):
                    world.grant_vehicle(account_id, f"veh-existing-{index}-{vehicle_index}", PARTS, operation_id=f"grant-existing-{index}-{vehicle_index}")
                if case["grant_allowed"]:
                    world.grant_vehicle(account_id, f"veh-candidate-{index}", PARTS, operation_id=f"grant-candidate-{index}")
                else:
                    with self.assertRaises(CapacityError):
                        world.grant_vehicle(account_id, f"veh-candidate-{index}", PARTS, operation_id=f"grant-candidate-{index}")

    def test_python_starter_deletion_matches_shared_vectors(self) -> None:
        for index, case in enumerate(self.vectors["routine_delete_cases"], start=1):
            with self.subTest(vector=case["name"]):
                world = WorldState()
                account_id = f"acct-delete-{index}"
                character_id = f"char-delete-{index}"
                vehicle_id = f"veh-delete-{index}"
                world.create_account(account_id, character_id, entitlements=Entitlements(garage_slots=2))
                world.grant_vehicle(account_id, vehicle_id, PARTS, starter_lineage=case["starter_lineage"], operation_id=f"grant-delete-{index}")
                if case["delete_allowed"]:
                    world.delete_vehicle(account_id, vehicle_id)
                    self.assertNotIn(vehicle_id, world.vehicles)
                else:
                    with self.assertRaises(DomainError):
                        world.delete_vehicle(account_id, vehicle_id)
                    self.assertIn(vehicle_id, world.vehicles)


if __name__ == "__main__":
    unittest.main()
