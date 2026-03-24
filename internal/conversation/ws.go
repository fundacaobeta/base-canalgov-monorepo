package conversation

import (
	"encoding/json"
	"time"

	cmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/conversation/models"
	wsmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/ws/models"
)

// BroadcastNewMessage broadcasts a new message to all users.
func (m *Manager) BroadcastNewMessage(message *cmodels.Message) {
	m.broadcastToUsers([]int{}, wsmodels.Message{
		Type: wsmodels.MessageTypeNewMessage,
		Data: map[string]interface{}{
			"conversation_uuid": message.ConversationUUID,
			"content":           "",
			"created_at":        message.CreatedAt.Format(time.RFC3339),
			"uuid":              message.UUID,
			"private":           message.Private,
			"type":              message.Type,
			"sender_type":       message.SenderType,
		},
	})
}

// BroadcastMessageUpdate broadcasts a message update to all users.
func (m *Manager) BroadcastMessageUpdate(conversationUUID, messageUUID, prop string, value any) {
	message := wsmodels.Message{
		Type: wsmodels.MessageTypeMessagePropUpdate,
		Data: map[string]interface{}{
			"conversation_uuid": conversationUUID,
			"uuid":              messageUUID,
			"prop":              prop,
			"value":             value,
		},
	}
	m.broadcastToUsers([]int{}, message)
}

// BroadcastConversationUpdate broadcasts a conversation update to all users.
func (m *Manager) BroadcastConversationUpdate(conversationUUID, prop string, value any) {
	message := wsmodels.Message{
		Type: wsmodels.MessageTypeConversationPropertyUpdate,
		Data: map[string]interface{}{
			"uuid":  conversationUUID,
			"prop":  prop,
			"value": value,
		},
	}
	m.broadcastToUsers([]int{}, message)
}

// BroadcastContactStatus broadcasts the online/offline status of a contact to all agents.
func (m *Manager) BroadcastContactStatus(contactID int, status string) {
	m.broadcastToUsers([]int{}, wsmodels.Message{
		Type: "contact_status_update",
		Data: map[string]interface{}{
			"contact_id": contactID,
			"status":     status,
		},
	})
}

// BroadcastTypingToConversation broadcasts a typing indicator to agents watching the conversation.
// If broadcastToWidgets is true, the typing indicator is also sent to widget clients connected to the conversation.
func (m *Manager) BroadcastTypingToConversation(conversationUUID string, isTyping bool, broadcastToWidgets bool) {
	m.broadcastToUsers([]int{}, wsmodels.Message{
		Type: "typing",
		Data: map[string]interface{}{
			"conversation_uuid":    conversationUUID,
			"is_typing":            isTyping,
			"broadcast_to_widgets": broadcastToWidgets,
		},
	})
}

// BroadcastContactPageVisit broadcasts a contact's page visit history to all agents.
func (m *Manager) BroadcastContactPageVisit(contactID int, pages []map[string]string) {
	m.broadcastToUsers([]int{}, wsmodels.Message{
		Type: "contact_page_visit",
		Data: map[string]interface{}{
			"contact_id": contactID,
			"pages":      pages,
		},
	})
}

// broadcastToUsers broadcasts a message to a list of users, if the list is empty it broadcasts to all users.
func (m *Manager) broadcastToUsers(userIDs []int, message wsmodels.Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		m.lo.Error("error marshalling WS message", "error", err)
		return
	}
	m.wsHub.BroadcastMessage(wsmodels.BroadcastMessage{
		Data:  messageBytes,
		Users: userIDs,
	})
}
