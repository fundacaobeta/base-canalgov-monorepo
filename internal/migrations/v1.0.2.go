package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_2 seeds notification settings for additional channels.
func V1_0_2(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		INSERT INTO settings ("key", value) VALUES
			('notification.whatsapp.enabled', 'false'::jsonb),
			('notification.whatsapp.provider', '"meta"'::jsonb),
			('notification.whatsapp.display_name', '""'::jsonb),
			('notification.whatsapp.phone_number_id', '""'::jsonb),
			('notification.whatsapp.access_token', '""'::jsonb),
			('notification.whatsapp.base_url', '""'::jsonb),
			('notification.whatsapp.webhook_verify_token', '""'::jsonb),
			('notification.whatsapp.default_template', '""'::jsonb),
			('notification.whatsapp.department_hint', '""'::jsonb),
			('notification.telegram.enabled', 'false'::jsonb),
			('notification.telegram.bot_name', '""'::jsonb),
			('notification.telegram.bot_token', '""'::jsonb),
			('notification.telegram.webhook_url', '""'::jsonb),
			('notification.telegram.default_chat_id', '""'::jsonb),
			('notification.telegram.allowed_updates', '"message,callback_query"'::jsonb),
			('notification.telegram.default_message', '""'::jsonb),
			('notification.sms.enabled', 'false'::jsonb),
			('notification.sms.provider', '"twilio"'::jsonb),
			('notification.sms.sender_id', '""'::jsonb),
			('notification.sms.api_key', '""'::jsonb),
			('notification.sms.api_secret', '""'::jsonb),
			('notification.sms.base_url', '""'::jsonb),
			('notification.sms.fallback_country_code', '"55"'::jsonb),
			('notification.sms.default_message', '""'::jsonb),
			('notification.push.enabled', 'false'::jsonb),
			('notification.push.provider', '"firebase"'::jsonb),
			('notification.push.app_id', '""'::jsonb),
			('notification.push.project_id', '""'::jsonb),
			('notification.push.api_key', '""'::jsonb),
			('notification.push.topic_default', '""'::jsonb),
			('notification.push.click_action_url', '""'::jsonb),
			('notification.push.payload_template', '"{\n  \"title\": \"Novo alerta\",\n  \"body\": \"Você recebeu uma nova atualização.\"\n}"'::jsonb),
			('notification.official_communications.enabled', 'false'::jsonb),
			('notification.official_communications.auto_create_conversation', 'true'::jsonb),
			('notification.official_communications.inbox_id', '""'::jsonb),
			('notification.official_communications.priority_id', '""'::jsonb),
			('notification.official_communications.status_id', '""'::jsonb),
			('notification.official_communications.subject_prefix', '"[Oficial] Comunicação recebida"'::jsonb),
			('notification.official_communications.target_sla_hours', '"24"'::jsonb),
			('notification.official_communications.default_types', '["Ofício", "Carta", "Notificação", "Intimação"]'::jsonb),
			('notification.official_communications.internal_note', '""'::jsonb),
			('notification.official_communications.routing_rules', '[]'::jsonb)
		ON CONFLICT ("key") DO NOTHING;
	`)
	return err
}
