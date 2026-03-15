-- Migration to add custom reports table
CREATE TABLE custom_reports (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    chart_type TEXT NOT NULL DEFAULT 'bar', -- bar, line, pie, metric
    metric_type TEXT NOT NULL DEFAULT 'conversations_count', -- conversations, response_time, sla_compliance
    filters JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_id INT REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_custom_reports_created_by ON custom_reports(created_by_id);
