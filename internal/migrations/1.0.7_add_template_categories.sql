-- Migration to add template categories and team associations
CREATE TABLE template_categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE template_category_teams (
    category_id INT REFERENCES template_categories(id) ON DELETE CASCADE,
    team_id INT REFERENCES teams(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, team_id)
);

ALTER TABLE templates ADD COLUMN category_id INT REFERENCES template_categories(id) ON DELETE SET NULL;

-- Add 'note' to template_type enum if it doesn't exist
-- Note: PostgreSQL doesn't support IF NOT EXISTS for ADD VALUE directly in some versions, 
-- but we can use a DO block.
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'note' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'template_type')) THEN
        ALTER TYPE template_type ADD VALUE 'note';
    END IF;
END $$;
