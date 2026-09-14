package core

// QuestSideEffects describes deterministic durable side effects derived only from
// a canonical MQ001..MQ100 quest identifier. It carries no player-supplied data.
type QuestSideEffects struct {
	InventoryItemID      string
	InventoryQuantity    int64
	BlueprintID          string
	MarkStarterRoadworthy bool
}

// QuestSideEffectsFor derives the authoritative quest side-effect schedule.
// PostgreSQL remains the durable authority that applies these effects transactionally.
func QuestSideEffectsFor(questID string) (QuestSideEffects, error) {
	if _, err := ParseQuestID(questID); err != nil {
		return QuestSideEffects{}, err
	}

	switch questID {
	case "MQ004":
		return QuestSideEffects{InventoryItemID: "part_brakes_track_i", InventoryQuantity: 1}, nil
	case "MQ005":
		return QuestSideEffects{BlueprintID: StarterRebuildBlueprint}, nil
	case "MQ009":
		return QuestSideEffects{InventoryItemID: "part_tires_street_i", InventoryQuantity: 1}, nil
	case "MQ012":
		return QuestSideEffects{MarkStarterRoadworthy: true}, nil
	default:
		return QuestSideEffects{}, nil
	}
}
