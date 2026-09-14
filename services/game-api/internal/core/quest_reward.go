package core

// QuestReward derives the canonical authoritative reward schedule for MQ001..MQ100.
// The caller supplies only the quest identifier; reward amounts are never client inputs.
func QuestReward(questID string) (RewardReceipt, error) {
	number, err := ParseQuestID(questID)
	if err != nil {
		return RewardReceipt{}, err
	}
	return RewardReceipt{
		QuestID:    questID,
		Applied:    true,
		Money:      int64(100 + number*10),
		XP:         int64(50 + number*5),
		Reputation: 1,
	}, nil
}
