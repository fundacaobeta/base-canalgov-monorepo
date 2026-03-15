package whatsapp

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
	ProviderWhatsApp = "whatsapp"
)

// WhatsApp implements the Notifier interface for sending messages via WhatsApp.
type WhatsApp struct {
	lo     *logf.Logger
	config models.WhatsAppNotification
	client *http.Client
}

// New initializes a new WhatsApp notifier.
func New(config models.WhatsAppNotification, lo *logf.Logger) (*WhatsApp, error) {
	return &WhatsApp{
		lo:     lo,
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send sends a notification message via WhatsApp API.
func (w *WhatsApp) Send(msg notifier.Message) error {
	if !w.config.Enabled {
		return fmt.Errorf("whatsapp notification is disabled")
	}

	// Basic implementation for Meta WhatsApp Cloud API
	// This would typically involve sending a template message or a session message.
	// For now, let's implement a generic structure for sending a text message.
	
	url := fmt.Sprintf("%s/%s/messages", w.config.BaseURL, w.config.PhoneNumberID)
	if w.config.BaseURL == "" {
		url = fmt.Sprintf("https://graph.facebook.com/v21.0/%s/messages", w.config.PhoneNumberID)
	}

	for _, recipient := range msg.RecipientEmails { // Reusing RecipientEmails as phone numbers for now
		w.lo.Info("sending whatsapp message", "recipient", recipient)
		
		// Note: Meta API requires a specific JSON payload. 
		// This is a simplified example.
		payload := fmt.Sprintf(`{
			"messaging_product": "whatsapp",
			"to": "%s",
			"type": "text",
			"text": {
				"body": "%s"
			}
		}`, recipient, strings.ReplaceAll(msg.Content, "\n", "\\n"))

		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+w.config.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := w.client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			w.lo.Error("whatsapp api error", "status", resp.Status)
			return fmt.Errorf("whatsapp api returned status: %s", resp.Status)
		}
	}

	return nil
}

// Name returns the name of the provider.
func (w *WhatsApp) Name() string {
	return ProviderWhatsApp
}
