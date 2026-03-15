// Package template manages templates including creation, retrieval and rendering.
package template
import (
	"database/sql"
	"embed"
	"errors"
	"html/template"
	"sync"


	"github.com/fundacaobeta/base-canalgov-monorepo/internal/dbutil"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/template/models"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/go-i18n"
	"github.com/zerodha/logf"
)

var (
	//go:embed queries.sql
	efs                   embed.FS
	ErrTemplateNotFound   = errors.New("template not found")
	TypeResponse          = "response"
	TypeEmailOutgoing     = "email_outgoing"
	TypeEmailNotification = "email_notification"
	TypeNote              = "note"
)

// Manager handles template-related operations.
type Manager struct {
	mutex   sync.RWMutex
	tpls    *template.Template
	webTpls *template.Template
	funcMap template.FuncMap
	q       queries
	lo      *logf.Logger
	i18n    *i18n.I18n
}

// queries contains prepared SQL queries.
type queries struct {
	InsertTemplate     *sqlx.Stmt `query:"insert"`
	UpdateTemplate     *sqlx.Stmt `query:"update"`
	DeleteTemplate     *sqlx.Stmt `query:"delete"`
	GetDefaultTemplate *sqlx.Stmt `query:"get-default"`
	GetAllTemplates    *sqlx.Stmt `query:"get-all"`
	GetAllByTeam       *sqlx.Stmt `query:"get-all-by-team"`
	GetTemplate        *sqlx.Stmt `query:"get-template"`
	GetByName          *sqlx.Stmt `query:"get-by-name"`
	IsBuiltIn          *sqlx.Stmt `query:"is-builtin"`

	// Category queries.
	GetCategories          *sqlx.Stmt `query:"get-categories"`
	GetCategoryTeams       *sqlx.Stmt `query:"get-category-teams"`
	InsertCategory         *sqlx.Stmt `query:"insert-category"`
	UpdateCategory         *sqlx.Stmt `query:"update-category"`
	DeleteCategory         *sqlx.Stmt `query:"delete-category"`
	ClearCategoryTeams     *sqlx.Stmt `query:"clear-category-teams"`
	InsertCategoryTeam     *sqlx.Stmt `query:"insert-category-team"`
	GetTemplatesByCategory *sqlx.Stmt `query:"get-templates-by-category"`
}

// New creates and returns a new instance of the Manager.
func New(lo *logf.Logger, db *sqlx.DB, webTpls *template.Template, tpls *template.Template, funcMap template.FuncMap, i18n *i18n.I18n) (*Manager, error) {
	var q queries
	if err := dbutil.ScanSQLFile("queries.sql", &q, db, efs); err != nil {
		return nil, err
	}
	return &Manager{
		mutex:   sync.RWMutex{},
		tpls:    tpls,
		webTpls: webTpls,
		funcMap: funcMap,
		q:       q,
		lo:      lo,
		i18n:    i18n,
	}, nil
}

// Update updates a new template with the given name, and body.
func (m *Manager) Update(id int, t models.Template) (models.Template, error) {
	var result models.Template
	if err := m.q.UpdateTemplate.Get(&result, id, t.Name, t.Body, t.IsDefault, t.Subject, t.Type, t.TeamID); err != nil {
		m.lo.Error("error updating template", "error", err)
		return models.Template{}, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.template}"), nil)
	}
	return result, nil
}

// Create creates a template.
func (m *Manager) Create(t models.Template) (models.Template, error) {
	if t.IsDefault && t.Type == "" {
		t.Type = TypeEmailOutgoing
	}
	var result models.Template
	if err := m.q.InsertTemplate.Get(&result, t.Name, t.Body, t.IsDefault, t.Subject, t.Type, t.TeamID); err != nil {
		if dbutil.IsUniqueViolationError(err) && t.IsDefault {
			return models.Template{}, envelope.NewError(envelope.GeneralError, m.i18n.T("template.defaultTemplateAlreadyExists"), nil)
		}
		m.lo.Error("error inserting template", "error", err)
		return models.Template{}, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorCreating", "name", "{globals.terms.template}"), nil)
	}
	return result, nil
}

// GetAll returns all templates by type.
func (m *Manager) GetAll(typ string, teamID *int, includeGlobal bool) ([]models.Template, error) {
	var templates = make([]models.Template, 0)
	var err error
	if teamID != nil {
		err = m.q.GetAllByTeam.Select(&templates, typ, *teamID, includeGlobal)
	} else {
		err = m.q.GetAllTemplates.Select(&templates, typ)
	}

	if err != nil {
		m.lo.Error("error fetching templates", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.template}"), nil)
	}
	return templates, nil
}

// Get returns a single template by ID.
func (m *Manager) Get(id int) (models.Template, error) {
	var t models.Template
	if err := m.q.GetTemplate.Get(&t, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrTemplateNotFound
		}
		m.lo.Error("error fetching template", "error", err)
		return t, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.template}"), nil)
	}
	return t, nil
}

// Delete deletes a template by ID.
func (m *Manager) Delete(id int) error {
	if _, err := m.q.DeleteTemplate.Exec(id); err != nil {
		m.lo.Error("error deleting template", "error", err)
		return envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorDeleting", "name", "{globals.terms.template}"), nil)
	}
	return nil
}

func (m *Manager) getDefaultOutgoingEmailTemplate() (models.Template, error) {
	var t models.Template
	if err := m.q.GetDefaultTemplate.Get(&t); err != nil {
		return t, err
	}
	return t, nil
}

func (m *Manager) getByName(name string) (models.Template, error) {
	var t models.Template
	if err := m.q.GetByName.Get(&t, name); err != nil {
		return t, err
	}
	return t, nil
}

// GetCategories returns all template categories.
func (m *Manager) GetCategories() ([]models.TemplateCategory, error) {
	var categories []models.TemplateCategory
	if err := m.q.GetCategories.Select(&categories); err != nil {
		m.lo.Error("error fetching template categories", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.category}"), nil)
	}

	for i := range categories {
		var teamIDs []int
		if err := m.q.GetCategoryTeams.Select(&teamIDs, categories[i].ID); err == nil {
			categories[i].TeamIDs = teamIDs
		}
	}

	return categories, nil
}

// CreateCategory creates a new template category.
func (m *Manager) CreateCategory(c models.TemplateCategory) (models.TemplateCategory, error) {
	var result models.TemplateCategory
	if err := m.q.InsertCategory.Get(&result, c.Name, c.Description); err != nil {
		m.lo.Error("error inserting template category", "error", err)
		return result, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorCreating", "name", "{globals.terms.category}"), nil)
	}

	for _, teamID := range c.TeamIDs {
		m.q.InsertCategoryTeam.Exec(result.ID, teamID)
	}
	result.TeamIDs = c.TeamIDs

	return result, nil
}

// UpdateCategory updates an existing template category.
func (m *Manager) UpdateCategory(id int, c models.TemplateCategory) (models.TemplateCategory, error) {
	var result models.TemplateCategory
	if err := m.q.UpdateCategory.Get(&result, id, c.Name, c.Description); err != nil {
		m.lo.Error("error updating template category", "error", err)
		return result, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.category}"), nil)
	}

	m.q.ClearCategoryTeams.Exec(id)
	for _, teamID := range c.TeamIDs {
		m.q.InsertCategoryTeam.Exec(id, teamID)
	}
	result.TeamIDs = c.TeamIDs

	return result, nil
}

// DeleteCategory deletes a template category.
func (m *Manager) DeleteCategory(id int) error {
	if _, err := m.q.DeleteCategory.Exec(id); err != nil {
		m.lo.Error("error deleting template category", "error", err)
		return envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorDeleting", "name", "{globals.terms.category}"), nil)
	}
	return nil
}

// Reload reloads the templates.
func (m *Manager) Reload(webTpls *template.Template, tpls *template.Template, funcMap template.FuncMap) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.webTpls = webTpls
	m.tpls = tpls
	m.funcMap = funcMap
	return nil
}
