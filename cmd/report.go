package main

import (
	"strconv"

	rmodels "github.com/fundacaobeta/base-canalgov-monorepo/internal/report/models"
	"github.com/zerodha/fastglue"
)

// handleOverviewCounts retrieves general dashboard counts for all users.
func handleOverviewCounts(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	counts, err := app.report.GetOverViewCounts()
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(counts)
}

// handleOverviewCharts retrieves general dashboard chart data.
func handleOverviewCharts(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		days, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("days")))
	)
	charts, err := app.report.GetOverviewChart(days)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(charts)
}

// handleOverviewSLA retrieves SLA data for the dashboard.
func handleOverviewSLA(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		days, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("days")))
	)
	sla, err := app.report.GetOverviewSLA(days)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(sla)
}

// handleOverviewCSAT retrieves CSAT metrics for the dashboard.
func handleOverviewCSAT(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		days, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("days")))
	)
	csat, err := app.report.GetOverviewCSAT(days)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(csat)
}

// handleOverviewMessageVolume retrieves message volume metrics for the dashboard.
func handleOverviewMessageVolume(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		days, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("days")))
	)
	volume, err := app.report.GetOverviewMessageVolume(days)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(volume)
}

// handleOverviewTagDistribution retrieves tag distribution metrics for the dashboard.
func handleOverviewTagDistribution(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		days, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("days")))
	)
	tags, err := app.report.GetOverviewTagDistribution(days)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(tags)
}

// handleGetCustomReports returns all custom reports.
func handleGetCustomReports(r *fastglue.Request) error {
	app := r.Context.(*App)
	reports, err := app.report.GetCustomReports()
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(reports)
}

// handleGetCustomReport returns a single custom report.
func handleGetCustomReport(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	report, err := app.report.GetCustomReport(id)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(report)
}

// handleCreateCustomReport creates a new custom report.
func handleCreateCustomReport(r *fastglue.Request) error {
	app := r.Context.(*App)
	var req rmodels.CustomReport
	if err := r.Decode(&req, "json"); err != nil {
		return sendErrorEnvelope(r, err)
	}
	report, err := app.report.CreateCustomReport(req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(report)
}

// handleUpdateCustomReport updates a custom report.
func handleUpdateCustomReport(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	var req rmodels.CustomReport
	if err := r.Decode(&req, "json"); err != nil {
		return sendErrorEnvelope(r, err)
	}
	report, err := app.report.UpdateCustomReport(id, req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(report)
}

// handleDeleteCustomReport deletes a custom report.
func handleDeleteCustomReport(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err := app.report.DeleteCustomReport(id); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

// handleExecuteCustomReport executes a custom report and returns the data.
func handleExecuteCustomReport(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	data, err := app.report.ExecuteCustomReport(id)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(data)
}
