package models

import (
	"time"

	"github.com/volatiletech/null/v9"
)

type TemplateType string

const (
	TemplateTypeResponse          TemplateType = "response"
	TemplateTypeEmailOutgoing     TemplateType = "email_outgoing"
	TemplateTypeEmailNotification TemplateType = "email_notification"
	TemplateTypeNote              TemplateType = "note"
)

type Template struct {
	ID         int         `db:"id" json:"id"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	Type       string      `db:"type" json:"type"`
	Name       string      `db:"name" json:"name"`
	Subject    null.String `db:"subject" json:"subject"`
	Body       string      `db:"body" json:"body"`
	IsDefault  bool        `db:"is_default" json:"is_default"`
	IsBuiltIn  bool        `db:"is_builtin" json:"is_builtin"`
	TeamID     *int        `db:"team_id" json:"team_id,omitempty"`
	TeamName   null.String `db:"team_name" json:"team_name"`
	CategoryID *int        `db:"category_id" json:"category_id,omitempty"`
}

type TemplateCategory struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	TeamIDs     []int     `db:"-" json:"team_ids,omitempty"`
}
