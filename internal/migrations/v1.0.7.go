package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V1_0_7 adds template categories and team associations.
func V1_0_7(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS template_categories (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS template_category_teams (
			category_id INT REFERENCES template_categories(id) ON DELETE CASCADE,
			team_id INT REFERENCES teams(id) ON DELETE CASCADE,
			PRIMARY KEY (category_id, team_id)
		);

		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='templates' AND column_name='category_id') THEN
				ALTER TABLE templates ADD COLUMN category_id INT REFERENCES template_categories(id) ON DELETE SET NULL;
			END IF;
		END $$;

		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'note' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'template_type')) THEN
				ALTER TYPE template_type ADD VALUE 'note';
			END IF;
		END $$;
	`)
	return err
}
