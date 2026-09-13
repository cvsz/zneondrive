package core

import "testing"

func TestNormalizeRaceID(t *testing.T) {
	if got, err := NormalizeRaceID(" race_first_ignition "); err != nil || got != "race_first_ignition" {
		t.Fatalf("unexpected normalized race: %q %v", got, err)
	}
	for _, invalid := range []string{"race", "Race_First", "race bad", "event_first"} {
		if _, err := NormalizeRaceID(invalid); err == nil {
			t.Fatalf("expected invalid race id: %q", invalid)
		}
	}
}

func TestRaceResultHashDeterministicAndBuildBound(t *testing.T) {
	base := RaceInstance{
		RaceInstanceID: "raceinst_test",
		RaceID: "race_first_ignition",
		AccountID: "acct_test",
		CharacterID: "char_test",
		VehicleID: "veh_test",
		BuildRevision: 4,
		BuildValidationHash: "buildhash-a",
	}
	a, err := RaceResultHash(base, 3, 120000)
	if err != nil {
		t.Fatal(err)
	}
	b, err := RaceResultHash(base, 3, 120000)
	if err != nil || a != b {
		t.Fatalf("hash should be deterministic: %q %q %v", a, b, err)
	}
	base.BuildValidationHash = "buildhash-b"
	c, err := RaceResultHash(base, 3, 120000)
	if err != nil {
		t.Fatal(err)
	}
	if a == c {
		t.Fatal("result hash must bind immutable build validation hash")
	}
}

func TestValidateRaceCheckpoint(t *testing.T) {
	if err := ValidateRaceCheckpoint(0, 1); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRaceCheckpoint(-1, 1); err == nil {
		t.Fatal("negative checkpoint must fail")
	}
	if err := ValidateRaceCheckpoint(0, 0); err == nil {
		t.Fatal("non-positive elapsed time must fail")
	}
}
