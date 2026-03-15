package push

import (
	"fmt"
	"net/http"
	"strings"

	notifier "github.com/fundacaobeta/base-canalgov-monorepo/internal/notification"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/setting/models"
	"github.com/zerodha/logf"
)

const (
	ProviderPush = "push"
)

// Push implements the Notifier interface for sending push notifications.
type Push struct {
	lo     *logf.Logger
	config models.PushNotification
	client *http.Client
}

// New initializes a new Push notifier.
func New(config models.PushNotification, lo *logf.Logger) (*Push, error) {
	return &Push{
		lo:     lo,
		config: config,
		client: &http.Client{},
	}, nil
}

// Send sends a push notification message.
func (p *Push) Send(msg notifier.Message) error {
	if !p.config.Enabled {
		return fmt.Errorf("push notification is disabled")
	}

	// Example for Firebase Cloud Messaging (FCM) logic
	url := "https://fcm.googleapis.com/fcm/send"
	if p.config.Provider == "onesignal" {
		url = "https://onesignal.com/api/v1/notifications"
	}

	for _, recipient := range msg.RecipientEmails { // Reusing as Device Tokens
		p.lo.Info("sending push notification", "token", recipient)

		payload := fmt.Sprintf(`{
			"to": "%s",
			"notification": {
				"title": "%s",
				"body": "%s"
			}
		}`, recipient, msg.Subject, msg.Content)

		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", "key="+p.config.APIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("push api returned status: %s", resp.Status)
		}
	}

	return nil
}

// Name returns the name of the provider.
func (p *Push) Name() string {
	return ProviderPush
}
