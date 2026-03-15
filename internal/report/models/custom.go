package models

import (
	"encoding/json"
	"time"
)

type CustomReport struct {
	ID          int             `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	ChartType   string          `json:"chart_type" db:"chart_type"`
	MetricType  string          `json:"metric_type" db:"metric_type"`
	Filters     json.RawMessage `json:"filters" db:"filters"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	CreatedByID *int            `json:"created_by_id" db:"created_by_id"`
}

type CustomReportResult struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}
