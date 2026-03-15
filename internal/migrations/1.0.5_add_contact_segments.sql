-- Migration to add contact segments table
CREATE TABLE contact_segments (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    filters JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for performance on filters if needed, though they are application-parsed
CREATE INDEX idx_contact_segments_name ON contact_segments(name);
