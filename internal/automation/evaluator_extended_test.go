package automation_test

// Test: Rule with multiple event triggers
func TestEventTrigger_MultipleEvents(t *testing.T) {
	mockStore := new(mockConversationStore)
	mockStore.On("ApplyAction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	engine := createTestEngine(mockStore)

	conversation := createTestConversation(func(c *cmodels.Conversation) {
		c.StatusID = null.IntFrom(3) // Resolved
	})

	rule := models.Rule{
		Events: []string{models.EventConversationStatusChange, models.EventConversationPriorityChange},
		Groups: []models.RuleGroup{
			{
				LogicalOp: models.OperatorOR,
				Rules: []models.RuleDetail{
					{Field: models.ConversationStatus, Operator: models.RuleOperatorEquals, Value: "3", FieldType: models.FieldTypeConversationField},
				},
			},
		},
		Actions: []models.RuleAction{
			{Type: models.ActionSendCSAT, Value: []string{"0"}},
		},
		GroupOperator: models.OperatorOR,
		ExecutionMode: models.ExecutionModeAll,
	}

	rules := []models.RuleRecord{
		{
			ID:      2,
			Name:    "Multi-Event Rule",
			Enabled: true,
			Events:  pq.StringArray{models.EventConversationStatusChange, models.EventConversationPriorityChange},
			Rules:   mustMarshal(rule),
		},
	}

	// Test with the first event
	engine.evalRules(rules, models.EventConversationStatusChange, conversation, umodels.User{})
	assert.Equal(t, 1, mockStore.callCount, "Should trigger on the first event")

	// Test with the second event
	engine.evalRules(rules, models.EventConversationPriorityChange, conversation, umodels.User{})
	assert.Equal(t, 2, mockStore.callCount, "Should trigger on the second event")

	// Test with an unrelated event
	engine.evalRules(rules, models.EventConversationMessageIncoming, conversation, umodels.User{})
	assert.Equal(t, 2, mockStore.callCount, "Should not trigger on an unrelated event")
}
