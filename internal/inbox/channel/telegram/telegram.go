package telegram

import (
	"context"

	cmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/conversation/models"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/inbox"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/inbox/models"
	"github.com/zerodha/logf"
)

// Telegram implements the Inbox interface for Telegram channel.
type Telegram struct {
	id       int
	lo       *logf.Logger
	msgStore inbox.MessageStore
	usrStore inbox.UserStore
	config   models.Config
}

// Opts contains options for creating a new Telegram inbox.
type Opts struct {
	ID       int
	Config   models.Config
	Lo       *logf.Logger
	MsgStore inbox.MessageStore
	UsrStore inbox.UserStore
}

// New initializes a new Telegram inbox.
func New(msgStore inbox.MessageStore, usrStore inbox.UserStore, opts Opts) (*Telegram, error) {
	return &Telegram{
		id:       opts.ID,
		lo:       opts.Lo,
		msgStore: msgStore,
		usrStore: usrStore,
		config:   opts.Config,
	}, nil
}

// Receive handles incoming messages via webhook.
func (t *Telegram) Receive(ctx context.Context) error {
	return nil
}

// Send sends a message to Telegram.
func (t *Telegram) Send(msg cmodels.OutboundMessage) error {
	t.lo.Info("sending message to telegram", "conversation_uuid", msg.ConversationUUID)
	return nil
}

// Close closes the Telegram inbox.
func (t *Telegram) Close() error {
	return nil
}

// Identifier returns the unique identifier for the inbox.
func (t *Telegram) Identifier() int {
	return t.id
}

// FromAddress returns the Telegram Bot name or ID.
func (t *Telegram) FromAddress() string {
	return t.config.From
}

// Channel returns the channel name.
func (t *Telegram) Channel() string {
	return inbox.ChannelTelegram
}

// HandleWebhook processes the incoming webhook payload from Telegram.
func (t *Telegram) HandleWebhook(payload []byte) error {
	t.lo.Debug("received telegram webhook payload", "payload", string(payload))
	return nil
}
