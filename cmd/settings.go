package main

import (
	"encoding/json"
	"net/mail"
	"strings"

	"github.com/abhinavxd/libredesk/internal/envelope"
	"github.com/abhinavxd/libredesk/internal/setting/models"
	"github.com/abhinavxd/libredesk/internal/stringutil"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

func maskSecret(value string) string {
	if value == "" {
		return ""
	}
	return strings.Repeat(stringutil.PasswordDummy, 10)
}

// handleGetGeneralSettings fetches general settings, this endpoint is not behind auth as it has no sensitive data and is required for the app to function.
func handleGetGeneralSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	out, err := app.setting.GetByPrefix("app")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	// Unmarshal to set the app.update to the settings, so the frontend can show that an update is available.
	var settings map[string]interface{}
	if err := json.Unmarshal(out, &settings); err != nil {
		app.lo.Error("error unmarshalling settings", "err", err)
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	// Set the app.update to the settings, adding `app` prefix to the key to match the settings structure in db.
	settings["app.update"] = app.update
	// Set app version.
	settings["app.version"] = versionString
	// Set restart required flag.
	settings["app.restart_required"] = app.restartRequired
	return r.SendEnvelope(settings)
}

// handleUpdateGeneralSettings updates general settings.
func handleUpdateGeneralSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.General{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	// Trim whitespace from string fields.
	req.SiteName = strings.TrimSpace(req.SiteName)
	req.FaviconURL = strings.TrimSpace(req.FaviconURL)
	req.LogoURL = strings.TrimSpace(req.LogoURL)
	req.Timezone = strings.TrimSpace(req.Timezone)
	// Trim whitespace and trailing slash from root URL.
	req.RootURL = strings.TrimRight(strings.TrimSpace(req.RootURL), "/")

	// Get current language before update.
	app.Lock()
	oldLang := ko.String("app.lang")
	app.Unlock()

	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	// Reload the settings and templates.
	if err := reloadSettings(app); err != nil {
		return envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.couldNotReload", "name", app.i18n.T("globals.terms.setting")), nil)
	}

	// Check if language changed and reload i18n if needed.
	app.Lock()
	newLang := ko.String("app.lang")
	if oldLang != newLang {
		app.lo.Info("language changed, reloading i18n", "old_lang", oldLang, "new_lang", newLang)
		app.i18n = initI18n(app.fs)
		app.lo.Info("reloaded i18n", "old_lang", oldLang, "new_lang", newLang)
	}
	app.Unlock()

	if err := reloadTemplates(app); err != nil {
		return envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.couldNotReload", "name", app.i18n.T("globals.terms.setting")), nil)
	}
	return r.SendEnvelope(true)
}

// handleGetEmailNotificationSettings fetches email notification settings.
func handleGetEmailNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.EmailNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.email")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}

	// Unmarshal and filter out password.
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	if notif.Password != "" {
		notif.Password = strings.Repeat(stringutil.PasswordDummy, 10)
	}
	return r.SendEnvelope(notif)
}

// handleUpdateEmailNotificationSettings updates email notification settings.
func handleUpdateEmailNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.EmailNotification{}
		cur = models.EmailNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	// Trim whitespace from string fields (Password intentionally NOT trimmed).
	req.Host = strings.TrimSpace(req.Host)
	req.Username = strings.TrimSpace(req.Username)
	req.EmailAddress = strings.TrimSpace(req.EmailAddress)
	req.HelloHostname = strings.TrimSpace(req.HelloHostname)
	req.IdleTimeout = strings.TrimSpace(req.IdleTimeout)
	req.WaitTimeout = strings.TrimSpace(req.WaitTimeout)

	out, err := app.setting.GetByPrefix("notification.email")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}

	if err := json.Unmarshal(out, &cur); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorUpdating", "name", app.i18n.T("globals.terms.setting")), nil))
	}

	// Make sure it's a valid from email address.
	if _, err := mail.ParseAddress(req.EmailAddress); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.invalidFromAddress"), nil, envelope.InputError)
	}

	// Retain current password if not changed.
	if req.Password == "" {
		req.Password = cur.Password
	}

	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}

	// Email notification settings require app restart to take effect.
	app.Lock()
	app.restartRequired = true
	app.Unlock()

	return r.SendEnvelope(true)
}

func handleGetWhatsAppNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.WhatsAppNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.whatsapp")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	notif.AccessToken = maskSecret(notif.AccessToken)
	notif.WebhookVerifyToken = maskSecret(notif.WebhookVerifyToken)
	return r.SendEnvelope(notif)
}

func handleUpdateWhatsAppNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.WhatsAppNotification{}
		cur = models.WhatsAppNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	req.Provider = strings.TrimSpace(req.Provider)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.PhoneNumberID = strings.TrimSpace(req.PhoneNumberID)
	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.DefaultTemplate = strings.TrimSpace(req.DefaultTemplate)
	req.DepartmentHint = strings.TrimSpace(req.DepartmentHint)

	out, err := app.setting.GetByPrefix("notification.whatsapp")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &cur); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorUpdating", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	if req.AccessToken == "" {
		req.AccessToken = cur.AccessToken
	}
	if req.WebhookVerifyToken == "" {
		req.WebhookVerifyToken = cur.WebhookVerifyToken
	}
	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleGetTelegramNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.TelegramNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.telegram")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	notif.BotToken = maskSecret(notif.BotToken)
	return r.SendEnvelope(notif)
}

func handleUpdateTelegramNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.TelegramNotification{}
		cur = models.TelegramNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	req.BotName = strings.TrimSpace(req.BotName)
	req.WebhookURL = strings.TrimSpace(req.WebhookURL)
	req.DefaultChatID = strings.TrimSpace(req.DefaultChatID)
	req.AllowedUpdates = strings.TrimSpace(req.AllowedUpdates)
	req.DefaultMessage = strings.TrimSpace(req.DefaultMessage)

	out, err := app.setting.GetByPrefix("notification.telegram")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &cur); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorUpdating", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	if req.BotToken == "" {
		req.BotToken = cur.BotToken
	}
	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleGetSMSNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.SMSNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.sms")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	notif.APIKey = maskSecret(notif.APIKey)
	notif.APISecret = maskSecret(notif.APISecret)
	return r.SendEnvelope(notif)
}

func handleUpdateSMSNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.SMSNotification{}
		cur = models.SMSNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	req.Provider = strings.TrimSpace(req.Provider)
	req.SenderID = strings.TrimSpace(req.SenderID)
	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.FallbackCountryCode = strings.TrimSpace(req.FallbackCountryCode)
	req.DefaultMessage = strings.TrimSpace(req.DefaultMessage)

	out, err := app.setting.GetByPrefix("notification.sms")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &cur); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorUpdating", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	if req.APIKey == "" {
		req.APIKey = cur.APIKey
	}
	if req.APISecret == "" {
		req.APISecret = cur.APISecret
	}
	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleGetPushNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.PushNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.push")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	notif.APIKey = maskSecret(notif.APIKey)
	return r.SendEnvelope(notif)
}

func handleUpdatePushNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.PushNotification{}
		cur = models.PushNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	req.Provider = strings.TrimSpace(req.Provider)
	req.AppID = strings.TrimSpace(req.AppID)
	req.ProjectID = strings.TrimSpace(req.ProjectID)
	req.TopicDefault = strings.TrimSpace(req.TopicDefault)
	req.ClickActionURL = strings.TrimSpace(req.ClickActionURL)
	req.PayloadTemplate = strings.TrimSpace(req.PayloadTemplate)

	out, err := app.setting.GetByPrefix("notification.push")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &cur); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorUpdating", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	if req.APIKey == "" {
		req.APIKey = cur.APIKey
	}
	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleGetOfficialCommunicationsNotificationSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		notif = models.OfficialCommunicationsNotification{}
	)

	out, err := app.setting.GetByPrefix("notification.official_communications")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &notif); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	return r.SendEnvelope(notif)
}

func handleGetMailDomainsSettings(r *fastglue.Request) error {
	var (
		app   = r.Context.(*App)
		mailD = models.MailDomains{}
	)

	out, err := app.setting.GetByPrefix("mail")
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	if err := json.Unmarshal(out, &mailD); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.GeneralError, app.i18n.Ts("globals.messages.errorFetching", "name", app.i18n.T("globals.terms.setting")), nil))
	}
	return r.SendEnvelope(mailD)
}

func handleUpdateMailDomainsSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.MailDomains{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	for i := range req.Domains {
		req.Domains[i].ID = strings.TrimSpace(req.Domains[i].ID)
		req.Domains[i].Name = strings.TrimSpace(req.Domains[i].Name)
		req.Domains[i].Domain = strings.TrimSpace(req.Domains[i].Domain)
		req.Domains[i].Provider = strings.TrimSpace(req.Domains[i].Provider)
		req.Domains[i].InboundStrategy = strings.TrimSpace(req.Domains[i].InboundStrategy)
		req.Domains[i].Notes = strings.TrimSpace(req.Domains[i].Notes)
	}

	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUpdateOfficialCommunicationsNotificationSettings(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.OfficialCommunicationsNotification{}
	)

	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.T("globals.messages.badRequest"), nil, envelope.InputError)
	}

	req.InboxID = strings.TrimSpace(req.InboxID)
	req.PriorityID = strings.TrimSpace(req.PriorityID)
	req.StatusID = strings.TrimSpace(req.StatusID)
	req.SubjectPrefix = strings.TrimSpace(req.SubjectPrefix)
	req.TargetSLAHours = strings.TrimSpace(req.TargetSLAHours)
	req.InternalNote = strings.TrimSpace(req.InternalNote)

	for i := range req.DefaultTypes {
		req.DefaultTypes[i] = strings.TrimSpace(req.DefaultTypes[i])
	}
	for i := range req.RoutingRules {
		req.RoutingRules[i].ID = strings.TrimSpace(req.RoutingRules[i].ID)
		req.RoutingRules[i].Name = strings.TrimSpace(req.RoutingRules[i].Name)
		for j := range req.RoutingRules[i].Types {
			req.RoutingRules[i].Types[j] = strings.TrimSpace(req.RoutingRules[i].Types[j])
		}
		for j := range req.RoutingRules[i].TeamIDs {
			req.RoutingRules[i].TeamIDs[j] = strings.TrimSpace(req.RoutingRules[i].TeamIDs[j])
		}
	}

	if err := app.setting.Update(req); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}
