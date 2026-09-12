package core

import "testing"

func TestBuildHashIsOrderIndependentAndDeduplicated(t *testing.T) {
	a, err := BuildHash([]string{"part_ecu_legacy_zero", "part_chassis_starter_prototype", "part_ecu_legacy_zero"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := BuildHash([]string{"part_chassis_starter_prototype", "part_ecu_legacy_zero"})
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("expected canonical build hash, got %s != %s", a, b)
	}
}

func TestQuestSequence(t *testing.T) {
	previous, ok, err := PreviousQuestID("MQ012")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || previous != "MQ011" {
		t.Fatalf("unexpected predecessor: %q %v", previous, ok)
	}
	if _, err := ParseQuestID("MQ101"); err == nil {
		t.Fatal("expected MQ101 to be rejected")
	}
}

func TestSecretGeneration(t *testing.T) {
	a, err := NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 64 || len(b) != 64 || a == b {
		t.Fatal("expected unique 256-bit hex secrets")
	}
}
