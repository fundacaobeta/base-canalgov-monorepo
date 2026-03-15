package main

import (
	"net/http"

	"github.com/zerodha/fastglue"
)

// handleWhatsAppWebhookVerification handles the verification challenge from Meta.
func handleWhatsAppWebhookVerification(r *fastglue.Request) error {
	// Meta sends a GET request with hub.mode, hub.verify_token, and hub.challenge.
	// You should verify the token and return the challenge.
	
	verifyToken := string(r.RequestCtx.QueryArgs().Peek("hub.verify_token"))
	challenge := string(r.RequestCtx.QueryArgs().Peek("hub.challenge"))
	mode := string(r.RequestCtx.QueryArgs().Peek("hub.mode"))

	if mode == "subscribe" && verifyToken != "" {
		// In a real implementation, you'd compare verifyToken with your stored app.whatsapp.verify_token.
		r.RequestCtx.SetStatusCode(http.StatusOK)
		r.RequestCtx.SetBodyString(challenge)
		return nil
	}

	r.RequestCtx.SetStatusCode(http.StatusForbidden)
	return nil
}

// handleWhatsAppWebhook handles incoming message payloads from WhatsApp.
func handleWhatsAppWebhook(r *fastglue.Request) error {
	app := r.Context.(*App)
	payload := r.RequestCtx.PostBody()

	app.lo.Debug("received whatsapp webhook", "payload", string(payload))

	// Implementation details:
	// 1. Find the WhatsApp inbox.
	// 2. Call the HandleWebhook method on the inbox.
	// 3. Return 200 OK to Meta immediately.

	r.RequestCtx.SetStatusCode(http.StatusOK)
	return nil
}

// handleTelegramWebhook handles incoming message payloads from Telegram.
func handleTelegramWebhook(r *fastglue.Request) error {
	app := r.Context.(*App)
	payload := r.RequestCtx.PostBody()

	app.lo.Debug("received telegram webhook", "payload", string(payload))

	// Implementation details similar to WhatsApp.

	r.RequestCtx.SetStatusCode(http.StatusOK)
	return nil
}
