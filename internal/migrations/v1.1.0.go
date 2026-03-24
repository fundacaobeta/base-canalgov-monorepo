package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_1_0 adds livechat and external user support:
// - external_user_id to users table
// - contact_last_seen_at, last_continuity_email_sent_at, last_interaction_sender_id to conversations
// - secret, linked_email_inbox_id to inboxes
// - Makes contact_channel_id nullable in conversations for livechat conversations
// - Adds livechat to channels enum
func V1_1_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		-- Add external_user_id to users
		ALTER TABLE users ADD COLUMN IF NOT EXISTS external_user_id TEXT NULL;
		CREATE INDEX IF NOT EXISTS idx_users_external_user_id ON users (external_user_id) WHERE external_user_id IS NOT NULL;

		-- Add livechat columns to conversations
		ALTER TABLE conversations ADD COLUMN IF NOT EXISTS contact_last_seen_at TIMESTAMPTZ NULL;
		ALTER TABLE conversations ADD COLUMN IF NOT EXISTS last_continuity_email_sent_at TIMESTAMPTZ NULL;
		ALTER TABLE conversations ADD COLUMN IF NOT EXISTS last_interaction_sender_id BIGINT REFERENCES users(id) ON DELETE SET NULL NULL;

		-- Make contact_channel_id nullable to support livechat conversations without a contact_channel
		ALTER TABLE conversations ALTER COLUMN contact_channel_id DROP NOT NULL;

		-- Add secret and linked_email_inbox_id to inboxes
		ALTER TABLE inboxes ADD COLUMN IF NOT EXISTS secret TEXT DEFAULT '' NULL;
		ALTER TABLE inboxes ADD COLUMN IF NOT EXISTS linked_email_inbox_id INT REFERENCES inboxes(id) ON DELETE SET NULL NULL;

		-- Add livechat to channels enum
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'livechat' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'channels')) THEN
				ALTER TYPE channels ADD VALUE 'livechat';
			END IF;
		END $$;
	`)
	return err
}
