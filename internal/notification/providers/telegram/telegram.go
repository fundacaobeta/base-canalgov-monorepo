package telegram

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	notifier "github.com/fundacaobeta/base-canalgov-monorepo/internal/notification"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/setting/models"
	"github.com/zerodha/logf"
)

const (
	ProviderTelegram = "telegram"
)

// Telegram implements the Notifier interface for sending messages via Telegram Bot API.
type Telegram struct {
	lo     *logf.Logger
	config models.TelegramNotification
	client *http.Client
}

// New initializes a new Telegram notifier.
func New(config models.TelegramNotification, lo *logf.Logger) (*Telegram, error) {
	return &Telegram{
		lo:     lo,
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send sends a notification message via Telegram Bot API.
func (t *Telegram) Send(msg notifier.Message) error {
	if !t.config.Enabled {
		return fmt.Errorf("telegram notification is disabled")
	}

	// Use the configured WebhookURL as base if provided, otherwise default to standard API
	apiBase := "https://api.telegram.org"
	url := fmt.Sprintf("%s/bot%s/sendMessage", apiBase, t.config.BotToken)

	for _, recipient := range msg.RecipientEmails { // Reusing RecipientEmails as ChatIDs for now
		t.lo.Info("sending telegram message", "chat_id", recipient)

		payload := fmt.Sprintf(`{
			"chat_id": "%s",
			"text": "%s",
			"parse_mode": "HTML"
		}`, recipient, strings.ReplaceAll(msg.Content, "\n", "\\n"))

		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := t.client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.lo.Error("telegram api error", "status", resp.Status)
			return fmt.Errorf("telegram api returned status: %s", resp.Status)
		}
	}

	return nil
}

// Name returns the name of the provider.
func (t *Telegram) Name() string {
	return ProviderTelegram
}
