package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_6 adds new channel values to the channels enum.
func V1_0_6(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'whatsapp' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'channels')) THEN
				ALTER TYPE channels ADD VALUE 'whatsapp';
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'telegram' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'channels')) THEN
				ALTER TYPE channels ADD VALUE 'telegram';
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'sms' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'channels')) THEN
				ALTER TYPE channels ADD VALUE 'sms';
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'push' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'channels')) THEN
				ALTER TYPE channels ADD VALUE 'push';
			END IF;
		END $$;
	`)
	return err
}
