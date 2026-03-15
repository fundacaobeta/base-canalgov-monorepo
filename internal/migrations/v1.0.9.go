package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_9 adds UNIQUE constraint to reset_password_token in users table.
func V1_0_9(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		ALTER TABLE users ADD CONSTRAINT users_reset_password_token_key UNIQUE (reset_password_token);
	`)
	return err
}
