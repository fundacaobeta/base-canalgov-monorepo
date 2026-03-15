package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_8 adds custom reports table.
func V1_0_8(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS custom_reports (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			chart_type TEXT NOT NULL DEFAULT 'bar',
			metric_type TEXT NOT NULL DEFAULT 'conversations_count',
			filters JSONB NOT NULL DEFAULT '[]'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_by_id INT REFERENCES users(id) ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_custom_reports_created_by ON custom_reports(created_by_id);
	`)
	return err
}
