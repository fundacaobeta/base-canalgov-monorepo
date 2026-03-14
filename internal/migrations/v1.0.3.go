package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_3 expands templates to support response models with optional team scope.
func V1_0_3(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	stmts := []string{
		`ALTER TYPE template_type ADD VALUE IF NOT EXISTS 'response'`,
		`ALTER TABLE templates ADD COLUMN IF NOT EXISTS team_id INT NULL REFERENCES teams(id) ON DELETE SET NULL ON UPDATE CASCADE`,
		`DROP INDEX IF EXISTS index_unique_templates_on_is_default_when_is_default_is_true`,
		`ALTER TABLE templates DROP CONSTRAINT IF EXISTS constraint_templates_on_team_for_response`,
		`ALTER TABLE templates ADD CONSTRAINT constraint_templates_on_team_for_response CHECK (
			team_id IS NULL OR type = 'response'
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS index_unique_templates_on_default_scope
			ON templates USING btree (type, COALESCE(team_id, 0))
			WHERE (is_default = true)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
