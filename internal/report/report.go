// Package report handles the management of reports.
package report

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

type customReportFilter struct {
	Field    string          `json:"field"`
	Operator string          `json:"operator"`
	Value    json.RawMessage `json:"value"`
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
	report, err := m.GetCustomReport(id)
	if err != nil {
		return nil, err
	}

	query, args, err := buildCustomReportQuery(report)
	if err != nil {
		return nil, err
	}

	var results []models.CustomReportResult
	if err := m.db.Select(&results, query, args...); err != nil {
		m.lo.Error("error executing custom report", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, m.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.report}"), nil)
	}

	return results, nil
}

func buildCustomReportQuery(report models.CustomReport) (string, []any, error) {
	selectClause := "COUNT(c.id)::float8 AS value"
	groupClause := ""
	orderClause := " ORDER BY value DESC, label ASC"
	joins := []string{}
	where := []string{"1=1"}
	args := []any{}

	switch report.MetricType {
	case "", "conversations_count":
		report.MetricType = "conversations_by_status"
		fallthrough
	case "conversations_by_status":
		selectClause = "COALESCE(cs.name, 'Sem status') AS label, COUNT(c.id)::float8 AS value"
		joins = append(joins, "LEFT JOIN conversation_statuses cs ON cs.id = c.status_id")
		groupClause = " GROUP BY cs.name"
	case "conversations_by_priority":
		selectClause = "COALESCE(cp.name, 'Sem prioridade') AS label, COUNT(c.id)::float8 AS value"
		joins = append(joins, "LEFT JOIN conversation_priorities cp ON cp.id = c.priority_id")
		groupClause = " GROUP BY cp.name"
	case "conversations_by_inbox":
		selectClause = "COALESCE(i.name, 'Sem caixa') AS label, COUNT(c.id)::float8 AS value"
		joins = append(joins, "LEFT JOIN inboxes i ON i.id = c.inbox_id")
		groupClause = " GROUP BY i.name"
	case "conversations_by_team":
		selectClause = "COALESCE(t.name, 'Sem equipe') AS label, COUNT(c.id)::float8 AS value"
		joins = append(joins, "LEFT JOIN teams t ON t.id = c.assigned_team_id")
		groupClause = " GROUP BY t.name"
	case "conversations_by_agent":
		selectClause = "COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), u.email, 'Sem agente') AS label, COUNT(c.id)::float8 AS value"
		joins = append(joins, "LEFT JOIN users u ON u.id = c.assigned_user_id")
		groupClause = " GROUP BY TRIM(CONCAT(u.first_name, ' ', u.last_name)), u.email"
	default:
		return "", nil, envelope.NewError(envelope.InputError, "Tipo de metrica de relatorio customizado invalido.", nil)
	}

	if len(report.Filters) > 0 && string(report.Filters) != "null" {
		filterClauses, filterArgs, err := buildCustomReportFilters(report.Filters, len(args)+1)
		if err != nil {
			return "", nil, err
		}
		where = append(where, filterClauses...)
		args = append(args, filterArgs...)
	}

	if report.ChartType == "metric" {
		selectClause = "'Total' AS label, COUNT(c.id)::float8 AS value"
		groupClause = ""
		orderClause = ""
	}

	query := fmt.Sprintf(
		"SELECT %s FROM conversations c %s WHERE %s%s%s",
		selectClause,
		strings.Join(joins, " "),
		strings.Join(where, " AND "),
		groupClause,
		orderClause,
	)

	return query, args, nil
}

func buildCustomReportFilters(raw json.RawMessage, startIndex int) ([]string, []any, error) {
	var filters []customReportFilter
	if err := json.Unmarshal(raw, &filters); err != nil {
		return nil, nil, envelope.NewError(envelope.InputError, "Filtros do relatorio customizado sao invalidos.", nil)
	}

	allowedFields := map[string]string{
		"status_id":        "c.status_id",
		"priority_id":      "c.priority_id",
		"assigned_team_id": "c.assigned_team_id",
		"assigned_user_id": "c.assigned_user_id",
		"inbox_id":         "c.inbox_id",
	}

	clauses := []string{}
	args := []any{}
	argIndex := startIndex

	for _, filter := range filters {
		if filter.Field == "" || filter.Operator == "" {
			continue
		}

		if filter.Field == "tags" {
			tagClauses, tagArgs, nextIndex, err := buildTagFilterClause(filter, argIndex)
			if err != nil {
				return nil, nil, err
			}
			if tagClauses != "" {
				clauses = append(clauses, tagClauses)
				args = append(args, tagArgs...)
				argIndex = nextIndex
			}
			continue
		}

		column, ok := allowedFields[filter.Field]
		if !ok {
			continue
		}

		clause, filterArgs, nextIndex, err := buildScalarFilterClause(column, filter, argIndex)
		if err != nil {
			return nil, nil, err
		}
		if clause == "" {
			continue
		}

		clauses = append(clauses, clause)
		args = append(args, filterArgs...)
		argIndex = nextIndex
	}

	return clauses, args, nil
}

func buildScalarFilterClause(column string, filter customReportFilter, argIndex int) (string, []any, int, error) {
	switch filter.Operator {
	case "set":
		return fmt.Sprintf("%s IS NOT NULL", column), nil, argIndex, nil
	case "not set":
		return fmt.Sprintf("%s IS NULL", column), nil, argIndex, nil
	}

	value, err := normalizeFilterValue(filter.Value)
	if err != nil {
		return "", nil, argIndex, err
	}

	switch filter.Operator {
	case "equals":
		return fmt.Sprintf("%s = $%d", column, argIndex), []any{value}, argIndex + 1, nil
	case "not equals":
		return fmt.Sprintf("%s <> $%d", column, argIndex), []any{value}, argIndex + 1, nil
	default:
		return "", nil, argIndex, nil
	}
}

func buildTagFilterClause(filter customReportFilter, argIndex int) (string, []any, int, error) {
	switch filter.Operator {
	case "set":
		return "EXISTS (SELECT 1 FROM conversation_tags ct WHERE ct.conversation_id = c.id)", nil, argIndex, nil
	case "not set":
		return "NOT EXISTS (SELECT 1 FROM conversation_tags ct WHERE ct.conversation_id = c.id)", nil, argIndex, nil
	}

	var rawValues []any
	if err := json.Unmarshal(filter.Value, &rawValues); err != nil {
		return "", nil, argIndex, envelope.NewError(envelope.InputError, "Filtro de tags invalido no relatorio customizado.", nil)
	}

	values := make([]any, 0, len(rawValues))
	placeholders := make([]string, 0, len(rawValues))
	for _, raw := range rawValues {
		switch v := raw.(type) {
		case string:
			if v == "" {
				continue
			}
			values = append(values, v)
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			argIndex++
		case float64:
			values = append(values, int(v))
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			argIndex++
		}
	}

	if len(placeholders) == 0 {
		return "", nil, argIndex, nil
	}

	subquery := fmt.Sprintf(
		"SELECT 1 FROM conversation_tags ct WHERE ct.conversation_id = c.id AND ct.tag_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	switch filter.Operator {
	case "contains":
		return fmt.Sprintf("EXISTS (%s)", subquery), values, argIndex, nil
	case "not contains":
		return fmt.Sprintf("NOT EXISTS (%s)", subquery), values, argIndex, nil
	default:
		return "", nil, argIndex, nil
	}
}

func normalizeFilterValue(raw json.RawMessage) (any, error) {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, envelope.NewError(envelope.InputError, "Valor de filtro invalido no relatorio customizado.", nil)
	}

	switch v := value.(type) {
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed, nil
		}
		return v, nil
	case float64:
		return int(v), nil
	default:
		return v, nil
	}
}
