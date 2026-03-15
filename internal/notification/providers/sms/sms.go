package sms

import (
	"fmt"
	"net/http"
	"strings"

	notifier "github.com/fundacaobeta/base-canalgov-monorepo/internal/notification"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/setting/models"
	"github.com/zerodha/logf"
)

const (
	ProviderSMS = "sms"
)

// SMS implements the Notifier interface for sending SMS via configured provider API.
type SMS struct {
	lo     *logf.Logger
	config models.SMSNotification
	client *http.Client
}

// New initializes a new SMS notifier.
func New(config models.SMSNotification, lo *logf.Logger) (*SMS, error) {
	return &SMS{
		lo:     lo,
		config: config,
		client: &http.Client{},
	}, nil
}

// Send sends an SMS message via provider API.
func (s *SMS) Send(msg notifier.Message) error {
	if !s.config.Enabled {
		return fmt.Errorf("sms notification is disabled")
	}

	for _, recipient := range msg.RecipientEmails { // Reusing RecipientEmails as phone numbers
		s.lo.Info("sending sms message", "to", recipient)
		
		url := s.config.BaseURL
		if url == "" {
			s.lo.Warn("sms base url is empty, cannot send")
			continue
		}

		payload := fmt.Sprintf(`{"to": "%s", "message": "%s"}`, recipient, msg.Content)
		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("sms api returned status: %s", resp.Status)
		}
	}

	return nil
}

// Name returns the name of the provider.
func (s *SMS) Name() string {
	return ProviderSMS
}
