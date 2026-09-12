package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrInvalidQuest = errors.New("invalid quest id")
	ErrInvalidParts = errors.New("invalid vehicle parts")
)

type InventoryItem struct {
	ItemID   string `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

type Snapshot struct {
	AccountID           string          `json:"account_id"`
	CharacterID         string          `json:"character_id"`
	VehicleID           string          `json:"vehicle_id"`
	Money               int64           `json:"money"`
	XP                  int64           `json:"xp"`
	Reputation          int64           `json:"reputation"`
	ActiveBuildRevision int             `json:"active_build_revision"`
	ActivePartIDs       []string        `json:"active_part_ids"`
	Inventory           []InventoryItem `json:"inventory"`
	Blueprints          []string        `json:"blueprints"`
	StarterLineage      bool            `json:"starter_lineage"`
	Roadworthy          bool            `json:"roadworthy"`
	CompletedQuests     []string        `json:"completed_quests"`
}

type RewardReceipt struct {
	QuestID    string `json:"quest_id"`
	Applied    bool   `json:"applied"`
	Money      int64  `json:"money"`
	XP         int64  `json:"xp"`
	Reputation int64  `json:"reputation"`
}

func NewSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func NewID(prefix string) (string, error) {
	raw, err := NewSecret()
	if err != nil {
		return "", err
	}
	return prefix + "_" + raw[:24], nil
}

func HashSecret(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func NormalizeParts(partIDs []string) ([]string, error) {
	if len(partIDs) == 0 || len(partIDs) > 32 {
		return nil, ErrInvalidParts
	}
	seen := make(map[string]struct{}, len(partIDs))
	normalized := make([]string, 0, len(partIDs))
	for _, partID := range partIDs {
		partID = strings.TrimSpace(partID)
		if !strings.HasPrefix(partID, "part_") || len(partID) > 128 || !IsCatalogPart(partID) {
			return nil, ErrInvalidParts
		}
		if _, exists := seen[partID]; exists {
			continue
		}
		seen[partID] = struct{}{}
		normalized = append(normalized, partID)
	}
	sort.Strings(normalized)
	if len(normalized) == 0 {
		return nil, ErrInvalidParts
	}
	return normalized, nil
}

func BuildHash(partIDs []string) (string, error) {
	normalized, err := NormalizeParts(partIDs)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("marshal parts: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func ParseQuestID(questID string) (int, error) {
	if len(questID) != 5 || !strings.HasPrefix(questID, "MQ") {
		return 0, ErrInvalidQuest
	}
	number, err := strconv.Atoi(questID[2:])
	if err != nil || number < 1 || number > 100 {
		return 0, ErrInvalidQuest
	}
	return number, nil
}

func PreviousQuestID(questID string) (string, bool, error) {
	number, err := ParseQuestID(questID)
	if err != nil {
		return "", false, err
	}
	if number == 1 {
		return "", false, nil
	}
	return fmt.Sprintf("MQ%03d", number-1), true, nil
}
