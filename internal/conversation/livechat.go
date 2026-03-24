package conversation

import (
	"database/sql"
	"encoding/json"
	"fmt"

	amodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/automation/models"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/conversation/models"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	mmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/media/models"
	wmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/webhook/models"
)

// SendAutoReply inserts an outgoing message from an AI assistant as an auto-reply in a conversation.
// It is used by the AI engine to deliver generated responses to the contact.
func (c *Manager) SendAutoReply(media []mmodels.Media, inboxID, senderID, contactID int, conversationUUID, content string, metaMap map[string]any) (models.Message, error) {
	metaJSON, err := json.Marshal(metaMap)
	if err != nil {
		return models.Message{}, fmt.Errorf("marshalling meta: %w", err)
	}

	message := models.Message{
		ConversationUUID: conversationUUID,
		SenderID:         senderID,
		Type:             models.MessageOutgoing,
		SenderType:       models.SenderTypeAgent,
		Status:           models.MessageStatusSent,
		Content:          content,
		ContentType:      models.ContentTypeHTML,
		Private:          false,
		Media:            media,
		Meta:             metaJSON,
	}

	if err := c.InsertMessage(&message); err != nil {
		return models.Message{}, err
	}
	return message, nil
}

// GetContactChatConversations returns all chat conversations for a given contact and inbox.
func (c *Manager) GetContactChatConversations(contactID, inboxID int) ([]models.ChatConversation, error) {
	var conversations = make([]models.ChatConversation, 0)
	if err := c.q.GetContactChatConversations.Select(&conversations, contactID, inboxID); err != nil {
		if err == sql.ErrNoRows {
			return conversations, nil
		}
		c.lo.Error("error fetching contact chat conversations", "contact_id", contactID, "inbox_id", inboxID, "error", err)
		return conversations, envelope.NewError(envelope.GeneralError, c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.conversation}"), nil)
	}
	return conversations, nil
}

// UpdateConversationContactLastSeen updates the contact_last_seen_at timestamp for a conversation.
func (c *Manager) UpdateConversationContactLastSeen(uuid string) error {
	if _, err := c.q.UpdateConversationContactLastSeen.Exec(uuid); err != nil {
		c.lo.Error("error updating contact last seen for conversation", "uuid", uuid, "error", err)
		return envelope.NewError(envelope.GeneralError, c.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.conversation}"), nil)
	}
	return nil
}

// ProcessIncomingMessageHooks runs post-message hooks for an incoming livechat message:
// - triggers automation rules
// - triggers webhooks
// - re-opens the conversation if it was not Open
func (c *Manager) ProcessIncomingMessageHooks(conversationUUID string, isNewConversation bool) error {
	conversation, err := c.GetConversation(0, conversationUUID, "")
	if err != nil {
		return fmt.Errorf("fetching conversation: %w", err)
	}

	if isNewConversation {
		c.webhookStore.TriggerEvent(wmodels.EventConversationCreated, conversation)
		c.automation.EvaluateNewConversationRules(conversation)
		return nil
	}

	// Re-open the conversation if it's not already Open.
	systemUser, err := c.userStore.GetSystemUser()
	if err != nil {
		c.lo.Error("error fetching system user for re-open", "error", err)
	} else {
		if err := c.ReOpenConversation(conversationUUID, systemUser); err != nil {
			c.lo.Error("error re-opening conversation", "conversation_uuid", conversationUUID, "error", err)
		}
	}

	// Trigger automations on incoming message event.
	c.automation.EvaluateConversationUpdateRules(conversation, amodels.EventConversationMessageIncoming)

	return nil
}

// ProcessIncomingLiveChatMessage inserts an incoming livechat message and runs broadcasting.
// Returns the inserted message.
func (c *Manager) ProcessIncomingLiveChatMessage(message models.Message) (models.Message, error) {
	// Upload any attachments first.
	if err := c.uploadMessageAttachments(&message); err != nil {
		c.lo.Error("error uploading livechat message attachments", "conversation_uuid", message.ConversationUUID, "error", err)
		return message, fmt.Errorf("uploading attachments: %w", err)
	}

	if err := c.InsertMessage(&message); err != nil {
		return message, err
	}

	// Trigger automation for incoming message.
	conversation, err := c.GetConversation(message.ConversationID, message.ConversationUUID, "")
	if err == nil {
		c.automation.EvaluateConversationUpdateRules(conversation, amodels.EventConversationMessageIncoming)
	}

	return message, nil
}

// BuildWidgetConversationResponse builds the full widget conversation response including messages.
// When includeMessages is true, the recent non-private messages are fetched and converted to ChatMessage format.
func (c *Manager) BuildWidgetConversationResponse(conversation models.Conversation, includeMessages bool) (WidgetConversationResponse, error) {
	var resp WidgetConversationResponse

	// Build ChatConversation.
	chatConv, err := c.buildChatConversation(conversation)
	if err != nil {
		return resp, fmt.Errorf("building chat conversation: %w", err)
	}
	resp.Conversation = chatConv

	if includeMessages {
		messages, _, err := c.GetConversationMessages(conversation.UUID, 1, maxMessagesPerPage, nil, []string{models.MessageIncoming, models.MessageOutgoing})
		if err != nil {
			c.lo.Error("error fetching conversation messages for widget", "conversation_uuid", conversation.UUID, "error", err)
			// Non-fatal: return empty messages.
			messages = []models.Message{}
		}

		chatMessages := make([]models.ChatMessage, 0, len(messages))
		for _, msg := range messages {
			// Skip private messages.
			if msg.Private {
				continue
			}
			// Update attachment URLs.
			for i := range msg.Attachments {
				if msg.Attachments[i].UUID != "" && msg.Attachments[i].URL == "" {
					msg.Attachments[i].URL = c.mediaStore.GetSignedURL(msg.Attachments[i].UUID)
				}
			}
			chatMessages = append(chatMessages, models.ChatMessage{
				UUID:             msg.UUID,
				Status:           msg.Status,
				ConversationUUID: msg.ConversationUUID,
				CreatedAt:        msg.CreatedAt,
				Content:          msg.Content,
				TextContent:      msg.TextContent,
				Author:           msg.Author,
				Attachments:      msg.Attachments,
				Meta:             msg.Meta,
			})
		}
		resp.Messages = chatMessages
	}

	return resp, nil
}

// buildChatConversation converts a Conversation to ChatConversation format.
func (c *Manager) buildChatConversation(conversation models.Conversation) (models.ChatConversation, error) {
	var chatConv models.ChatConversation

	if err := c.q.GetChatConversation.Get(&chatConv, conversation.UUID); err != nil {
		if err == sql.ErrNoRows {
			return chatConv, envelope.NewError(envelope.NotFoundError, c.i18n.Ts("globals.messages.notFound", "name", "{globals.terms.conversation}"), nil)
		}
		c.lo.Error("error fetching chat conversation", "uuid", conversation.UUID, "error", err)
		return chatConv, envelope.NewError(envelope.GeneralError, c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.conversation}"), nil)
	}

	return chatConv, nil
}
