package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/conversation/models"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/inbox"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/inbox/models"
	"github.com/zerodha/logf"
)

// WhatsApp implements the Inbox interface for WhatsApp channel.
type WhatsApp struct {
	id       int
	lo       *logf.Logger
	msgStore inbox.MessageStore
	usrStore inbox.UserStore
	config   models.Config
}

// Opts contains options for creating a new WhatsApp inbox.
type Opts struct {
	ID       int
	Config   models.Config
	Lo       *logf.Logger
	MsgStore inbox.MessageStore
	UsrStore inbox.UserStore
}

// New initializes a new WhatsApp inbox.
func New(msgStore inbox.MessageStore, usrStore inbox.UserStore, opts Opts) (*WhatsApp, error) {
	return &WhatsApp{
		id:       opts.ID,
		lo:       opts.Lo,
		msgStore: msgStore,
		usrStore: usrStore,
		config:   opts.Config,
	}, nil
}

// Receive handles incoming messages via webhook.
// For WhatsApp, this is usually called by a handler in cmd/handlers.go that routes to this inbox.
func (w *WhatsApp) Receive(ctx context.Context) error {
	// Webhook-based channels don't need a background receiver loop like IMAP.
	// The actual receiving happens in the HTTP handler.
	return nil
}

// Send sends a message to WhatsApp.
func (w *WhatsApp) Send(msg cmodels.Message) error {
	// Implementation for sending a message back to the user via WhatsApp API.
	// This would use the Meta API similar to the notifier provider.
	w.lo.Info("sending message to whatsapp", "conversation", msg.ConversationID)
	return nil
}

// Close closes the WhatsApp inbox.
func (w *WhatsApp) Close() error {
	return nil
}

// Identifier returns the unique identifier for the inbox.
func (w *WhatsApp) Identifier() int {
	return w.id
}

// FromAddress returns the WhatsApp phone number or ID.
func (w *WhatsApp) FromAddress() string {
	return w.config.From
}

// Channel returns the channel name.
func (w *WhatsApp) Channel() string {
	return inbox.ChannelWhatsApp
}

// HandleWebhook processes the incoming webhook payload from Meta.
func (w *WhatsApp) HandleWebhook(payload []byte) error {
	// This is where the complex parsing of Meta's JSON occurs.
	w.lo.Debug("received whatsapp webhook payload", "payload", string(payload))
	
	// Implementation would:
	// 1. Parse payload
	// 2. Identify sender (phone number)
	// 3. Create/Find contact via w.usrStore
	// 4. Enqueue incoming message via w.msgStore.EnqueueIncoming
	
	return nil
}
