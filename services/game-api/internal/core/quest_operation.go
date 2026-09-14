package core

import "strings"

const questOperationScopeDomain = "zneondrive:quest-operation:v1"

// ScopeQuestOperationID derives the durable operation key stored for a quest
// mutation from authoritative character identity plus the caller's opaque
// idempotency key. The derivation prevents one character from replaying another
// character's quest operation while preserving replay across renewed sessions.
func ScopeQuestOperationID(characterID, operationID string) (string, bool) {
	characterID = strings.TrimSpace(characterID)
	operationID = strings.TrimSpace(operationID)
	if characterID == "" || operationID == "" {
		return "", false
	}
	return HashSecret(questOperationScopeDomain + "\x00" + characterID + "\x00" + operationID), true
}
