package core

import "testing"

func TestScopeQuestOperationID(t *testing.T) {
	first, ok := ScopeQuestOperationID("char_a", "op-123")
	if !ok || first == "" {
		t.Fatal("expected scoped quest operation id")
	}
	replay, ok := ScopeQuestOperationID("char_a", "op-123")
	if !ok || replay != first {
		t.Fatal("same character and operation must derive the same durable key")
	}
	otherCharacter, ok := ScopeQuestOperationID("char_b", "op-123")
	if !ok || otherCharacter == first {
		t.Fatal("different characters must not share a durable operation key")
	}
	otherOperation, ok := ScopeQuestOperationID("char_a", "op-456")
	if !ok || otherOperation == first {
		t.Fatal("different caller operation ids must derive distinct durable keys")
	}
	if _, ok := ScopeQuestOperationID("", "op-123"); ok {
		t.Fatal("empty character id must be rejected")
	}
	if _, ok := ScopeQuestOperationID("char_a", "   "); ok {
		t.Fatal("empty operation id must be rejected")
	}
}
