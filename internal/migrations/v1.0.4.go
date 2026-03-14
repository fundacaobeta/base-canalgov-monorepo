package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_4 seeds managed mail domain settings.
func V1_0_4(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		INSERT INTO settings ("key", value) VALUES
			('mail.domains', '[]'::jsonb)
		ON CONFLICT ("key") DO NOTHING;
	`)
	return err
}
