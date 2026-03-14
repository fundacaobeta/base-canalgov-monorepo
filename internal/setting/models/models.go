package models

type General struct {
	SiteName                    string   `json:"app.site_name"`
	Lang                        string   `json:"app.lang"`
	MaxFileUploadSize           int      `json:"app.max_file_upload_size"`
	FaviconURL                  string   `json:"app.favicon_url"`
	LogoURL                     string   `json:"app.logo_url"`
	RootURL                     string   `json:"app.root_url"`
	AllowedFileUploadExtensions []string `json:"app.allowed_file_upload_extensions"`
	Timezone                    string   `json:"app.timezone"`
	BusinessHoursID             string   `json:"app.business_hours_id"`
}

type MailDomain struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Domain          string `json:"domain"`
	Provider        string `json:"provider"`
	InboundStrategy string `json:"inbound_strategy"`
	Enabled         bool   `json:"enabled"`
	IsDefault       bool   `json:"is_default"`
	Notes           string `json:"notes"`
}

type MailDomains struct {
	Domains []MailDomain `json:"mail.domains" db:"mail.domains"`
}

type EmailNotification struct {
	Username      string `json:"notification.email.username" db:"notification.email.username"`
	Host          string `json:"notification.email.host" db:"notification.email.host"`
	Port          int    `json:"notification.email.port" db:"notification.email.port"`
	Password      string `json:"notification.email.password" db:"notification.email.password"`
	MaxConns      int    `json:"notification.email.max_conns" db:"notification.email.max_conns"`
	IdleTimeout   string `json:"notification.email.idle_timeout" db:"notification.email.idle_timeout"`
	WaitTimeout   string `json:"notification.email.wait_timeout" db:"notification.email.wait_timeout"`
	AuthProtocol  string `json:"notification.email.auth_protocol" db:"notification.email.auth_protocol"`
	EmailAddress  string `json:"notification.email.email_address" db:"notification.email.email_address"`
	MaxMsgRetries int    `json:"notification.email.max_msg_retries" db:"notification.email.max_msg_retries"`
	TLSType       string `json:"notification.email.tls_type" db:"notification.email.tls_type"`
	TLSSkipVerify bool   `json:"notification.email.tls_skip_verify" db:"notification.email.tls_skip_verify"`
	HelloHostname string `json:"notification.email.hello_hostname" db:"notification.email.hello_hostname"`
	Enabled       bool   `json:"notification.email.enabled" db:"notification.email.enabled"`
}

type WhatsAppNotification struct {
	Enabled            bool   `json:"notification.whatsapp.enabled" db:"notification.whatsapp.enabled"`
	Provider           string `json:"notification.whatsapp.provider" db:"notification.whatsapp.provider"`
	DisplayName        string `json:"notification.whatsapp.display_name" db:"notification.whatsapp.display_name"`
	PhoneNumberID      string `json:"notification.whatsapp.phone_number_id" db:"notification.whatsapp.phone_number_id"`
	AccessToken        string `json:"notification.whatsapp.access_token" db:"notification.whatsapp.access_token"`
	BaseURL            string `json:"notification.whatsapp.base_url" db:"notification.whatsapp.base_url"`
	WebhookVerifyToken string `json:"notification.whatsapp.webhook_verify_token" db:"notification.whatsapp.webhook_verify_token"`
	DefaultTemplate    string `json:"notification.whatsapp.default_template" db:"notification.whatsapp.default_template"`
	DepartmentHint     string `json:"notification.whatsapp.department_hint" db:"notification.whatsapp.department_hint"`
}

type TelegramNotification struct {
	Enabled        bool   `json:"notification.telegram.enabled" db:"notification.telegram.enabled"`
	BotName        string `json:"notification.telegram.bot_name" db:"notification.telegram.bot_name"`
	BotToken       string `json:"notification.telegram.bot_token" db:"notification.telegram.bot_token"`
	WebhookURL     string `json:"notification.telegram.webhook_url" db:"notification.telegram.webhook_url"`
	DefaultChatID  string `json:"notification.telegram.default_chat_id" db:"notification.telegram.default_chat_id"`
	AllowedUpdates string `json:"notification.telegram.allowed_updates" db:"notification.telegram.allowed_updates"`
	DefaultMessage string `json:"notification.telegram.default_message" db:"notification.telegram.default_message"`
}

type SMSNotification struct {
	Enabled             bool   `json:"notification.sms.enabled" db:"notification.sms.enabled"`
	Provider            string `json:"notification.sms.provider" db:"notification.sms.provider"`
	SenderID            string `json:"notification.sms.sender_id" db:"notification.sms.sender_id"`
	APIKey              string `json:"notification.sms.api_key" db:"notification.sms.api_key"`
	APISecret           string `json:"notification.sms.api_secret" db:"notification.sms.api_secret"`
	BaseURL             string `json:"notification.sms.base_url" db:"notification.sms.base_url"`
	FallbackCountryCode string `json:"notification.sms.fallback_country_code" db:"notification.sms.fallback_country_code"`
	DefaultMessage      string `json:"notification.sms.default_message" db:"notification.sms.default_message"`
}

type PushNotification struct {
	Enabled         bool   `json:"notification.push.enabled" db:"notification.push.enabled"`
	Provider        string `json:"notification.push.provider" db:"notification.push.provider"`
	AppID           string `json:"notification.push.app_id" db:"notification.push.app_id"`
	ProjectID       string `json:"notification.push.project_id" db:"notification.push.project_id"`
	APIKey          string `json:"notification.push.api_key" db:"notification.push.api_key"`
	TopicDefault    string `json:"notification.push.topic_default" db:"notification.push.topic_default"`
	ClickActionURL  string `json:"notification.push.click_action_url" db:"notification.push.click_action_url"`
	PayloadTemplate string `json:"notification.push.payload_template" db:"notification.push.payload_template"`
}

type OfficialCommunicationRoutingRule struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	TeamIDs []string `json:"team_ids"`
}

type OfficialCommunicationsNotification struct {
	Enabled                bool                               `json:"notification.official_communications.enabled" db:"notification.official_communications.enabled"`
	AutoCreateConversation bool                               `json:"notification.official_communications.auto_create_conversation" db:"notification.official_communications.auto_create_conversation"`
	InboxID                string                             `json:"notification.official_communications.inbox_id" db:"notification.official_communications.inbox_id"`
	PriorityID             string                             `json:"notification.official_communications.priority_id" db:"notification.official_communications.priority_id"`
	StatusID               string                             `json:"notification.official_communications.status_id" db:"notification.official_communications.status_id"`
	SubjectPrefix          string                             `json:"notification.official_communications.subject_prefix" db:"notification.official_communications.subject_prefix"`
	TargetSLAHours         string                             `json:"notification.official_communications.target_sla_hours" db:"notification.official_communications.target_sla_hours"`
	DefaultTypes           []string                           `json:"notification.official_communications.default_types" db:"notification.official_communications.default_types"`
	InternalNote           string                             `json:"notification.official_communications.internal_note" db:"notification.official_communications.internal_note"`
	RoutingRules           []OfficialCommunicationRoutingRule `json:"notification.official_communications.routing_rules" db:"notification.official_communications.routing_rules"`
}

type Settings struct {
	EmailNotification
	WhatsAppNotification
	TelegramNotification
	SMSNotification
	PushNotification
	OfficialCommunicationsNotification
	MailDomains
	General
}
