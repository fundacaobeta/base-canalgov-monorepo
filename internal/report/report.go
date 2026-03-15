// Package report handles the management of reports.
package report

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/dbutil"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/report/models"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/go-i18n"
	"github.com/zerodha/logf"
)

var (
	//go:embed queries.sql
	efs embed.FS
)

type Manager struct {
	q    queries
	lo   *logf.Logger
	i18n *i18n.I18n
	db   *sqlx.DB
}

// Opts contains options for initializing the report Manager.
type Opts struct {
	DB   *sqlx.DB
	Lo   *logf.Logger
	I18n *i18n.I18n
}

// queries contains prepared SQL queries.
type queries struct {
	GetOverviewCharts          string `query:"get-overview-charts"`
	GetOverviewCounts          string `query:"get-overview-counts"`
	GetOverviewSLA             string `query:"get-overview-sla-counts"`
	GetOverviewCSAT            string `query:"get-overview-csat"`
	GetOverviewMessageVolume   string `query:"get-overview-message-volume"`
	GetOverviewTagDistribution string `query:"get-overview-tag-distribution"`

	// Custom reports.
	GetCustomReports   *sqlx.Stmt `query:"get-custom-reports"`
	GetCustomReport    *sqlx.Stmt `query:"get-custom-report"`
	InsertCustomReport *sqlx.Stmt `query:"insert-custom-report"`
	UpdateCustomReport *sqlx.Stmt `query:"update-custom-report"`
	DeleteCustomReport *sqlx.Stmt `query:"delete-custom-report"`
}

// New creates and returns a new instance of the Manager.
func New(opts Opts) (*Manager, error) {
	var q queries
	if err := dbutil.ScanSQLFile("queries.sql", &q, opts.DB, efs); err != nil {
		return nil, err
	}
	return &Manager{
		q:    q,
		lo:   opts.Lo,
		i18n: opts.I18n,
		db:   opts.DB,
	}, nil
}

// GetOverViewCounts returns overview counts
func (m *Manager) GetOverViewCounts() (json.RawMessage, error) {
	var counts = json.RawMessage{}
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}
	defer tx.Rollback()

	if err := tx.Get(&counts, m.q.GetOverviewCounts); err != nil {
		m.lo.Error("error fetching overview counts", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}

	if err := tx.Commit(); err != nil {
		m.lo.Error("error committing db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}

	return counts, nil
}

// GetOverviewSLA returns overview SLA data
func (m *Manager) GetOverviewSLA(days int) (json.RawMessage, error) {
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}
	defer tx.Rollback()

	var result models.OverviewSLA
	// Format query with days parameter for both CTEs
	query := fmt.Sprintf(m.q.GetOverviewSLA, days, days, days, days)
	if err := tx.Get(&result, query); err != nil {
		m.lo.Error("error fetching overview SLA data", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}

	if err := tx.Commit(); err != nil {
		m.lo.Error("error committing db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}

	slaData, err := json.Marshal(result)
	if err != nil {
		m.lo.Error("error marshaling SLA data", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingCount", "name", "{globals.terms.overview}"), nil)
	}

	return slaData, nil
}

// GetOverviewChart returns overview chart data
func (m *Manager) GetOverviewChart(days int) (json.RawMessage, error) {
	var stats = json.RawMessage{}
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingChart", "name", "{globals.terms.overview}"), nil)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(m.q.GetOverviewCharts, days, days, days, days)
	if err := tx.Get(&stats, query); err != nil {
		m.lo.Error("error fetching overview charts", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetchingChart", "name", "{globals.terms.overview}"), nil)
	}
	return stats, nil
}

// GetOverviewCSAT returns CSAT metrics for the overview dashboard
func (m *Manager) GetOverviewCSAT(days int) (json.RawMessage, error) {
	var stats = json.RawMessage{}
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.csat}"), nil)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(m.q.GetOverviewCSAT, days, days)
	if err := tx.Get(&stats, query); err != nil {
		m.lo.Error("error fetching overview CSAT", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.csat}"), nil)
	}
	return stats, nil
}

// GetOverviewMessageVolume returns message volume metrics for the overview dashboard
func (m *Manager) GetOverviewMessageVolume(days int) (json.RawMessage, error) {
	var stats = json.RawMessage{}
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.message}"), nil)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(m.q.GetOverviewMessageVolume, days, days)
	if err := tx.Get(&stats, query); err != nil {
		m.lo.Error("error fetching overview message volume", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.message}"), nil)
	}
	return stats, nil
}

// GetOverviewTagDistribution returns tag distribution metrics for the overview dashboard
func (m *Manager) GetOverviewTagDistribution(days int) (json.RawMessage, error) {
	var stats = json.RawMessage{}
	tx, err := m.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		m.lo.Error("error starting db txn", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.tag}"), nil)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(m.q.GetOverviewTagDistribution, days, days, days, days)
	if err := tx.Get(&stats, query); err != nil {
		m.lo.Error("error fetching overview tag distribution", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.tag}"), nil)
	}
	return stats, nil
}

// GetCustomReports returns all custom reports.
func (m *Manager) GetCustomReports() ([]models.CustomReport, error) {
	var reports []models.CustomReport
	if err := m.q.GetCustomReports.Select(&reports); err != nil {
		m.lo.Error("error fetching custom reports", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.report}"), nil)
	}
	return reports, nil
}

// GetCustomReport returns a single custom report by ID.
func (m *Manager) GetCustomReport(id int) (models.CustomReport, error) {
	var report models.CustomReport
	if err := m.q.GetCustomReport.Get(&report, id); err != nil {
		m.lo.Error("error fetching custom report", "error", err)
		return report, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.report}"), nil)
	}
	return report, nil
}

// CreateCustomReport creates a new custom report.
func (m *Manager) CreateCustomReport(r models.CustomReport) (models.CustomReport, error) {
	var result models.CustomReport
	if err := m.q.InsertCustomReport.Get(&result, r.Name, r.Description, r.ChartType, r.MetricType, r.Filters, r.CreatedByID); err != nil {
		m.lo.Error("error creating custom report", "error", err)
		return result, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorCreating", "name", "{globals.terms.report}"), nil)
	}
	return result, nil
}

// UpdateCustomReport updates an existing custom report.
func (m *Manager) UpdateCustomReport(id int, r models.CustomReport) (models.CustomReport, error) {
	var result models.CustomReport
	if err := m.q.UpdateCustomReport.Get(&result, id, r.Name, r.Description, r.ChartType, r.MetricType, r.Filters); err != nil {
		m.lo.Error("error updating custom report", "error", err)
		return result, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.report}"), nil)
	}
	return result, nil
}

// DeleteCustomReport deletes a custom report.
func (m *Manager) DeleteCustomReport(id int) error {
	if _, err := m.q.DeleteCustomReport.Exec(id); err != nil {
		m.lo.Error("error deleting custom report", "error", err)
		return envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorDeleting", "name", "{globals.terms.report}"), nil)
	}
	return nil
}

// ExecuteCustomReport runs the custom report and returns the aggregated data.
func (m *Manager) ExecuteCustomReport(id int) ([]models.CustomReportResult, error) {
	_, err := m.GetCustomReport(id)
	if err != nil {
		return nil, err
	}

	// Dynamic aggregation based on filters.
	// For now, grouping by status as a baseline.
	query := `
		SELECT cs.name AS label, COUNT(c.id)::float8 AS value
		FROM conversations c
		JOIN conversation_statuses cs ON cs.id = c.status_id
		WHERE c.deleted_at IS NULL
		GROUP BY cs.name`

	var results []models.CustomReportResult
	if err := m.db.Select(&results, query); err != nil {
		m.lo.Error("error executing custom report", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.report}"), nil)
	}

	return results, nil
}
