package main

import (
	"strconv"

	"github.com/abhinavxd/libredesk/internal/envelope"
	"github.com/abhinavxd/libredesk/internal/template/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

// handleGetTemplates returns all templates.
func handleGetTemplates(r *fastglue.Request) error {
	var (
		app           = r.Context.(*App)
		typ           = string(r.RequestCtx.QueryArgs().Peek("type"))
		teamIDRaw     = string(r.RequestCtx.QueryArgs().Peek("team_id"))
		includeGlobal = string(r.RequestCtx.QueryArgs().Peek("include_global")) == "true"
		teamID        *int
	)
	if typ == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.empty", "name", "`type`"), nil, envelope.InputError)
	}
	if teamIDRaw != "" {
		id, err := strconv.Atoi(teamIDRaw)
		if err != nil || id <= 0 {
			return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.invalid", "name", "`team_id`"), nil, envelope.InputError)
		}
		teamID = &id
	}
	t, err := app.tmpl.GetAll(typ, teamID, includeGlobal)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(t)
}

// handleGetTemplate returns a template by id.
func handleGetTemplate(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	id, err := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err != nil || id == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.invalid", "name", "`id`"), nil, envelope.InputError)
	}
	t, err := app.tmpl.Get(id)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(t)
}

// handleCreateTemplate creates a new template.
func handleCreateTemplate(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.Template{}
	)
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.errorParsing", "name", "{globals.terms.request}"), nil, envelope.InputError)
	}
	if req.Name == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.empty", "name", "`name`"), nil, envelope.InputError)
	}
	if req.Type == "" {
		req.Type = "response"
	}
	if req.Type != "response" {
		req.TeamID = nil
	}
	template, err := app.tmpl.Create(req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(template)
}

// handleUpdateTemplate updates a template.
func handleUpdateTemplate(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.Template{}
	)
	id, err := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err != nil || id == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Invalid template `id`.", nil, envelope.InputError)
	}
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.errorParsing", "name", "{globals.terms.request}"), nil, envelope.InputError)
	}
	if req.Name == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.empty", "name", "`name`"), nil, envelope.InputError)
	}
	if req.Type == "" {
		req.Type = "response"
	}
	if req.Type != "response" {
		req.TeamID = nil
	}
	updatedTemplate, err := app.tmpl.Update(id, req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(updatedTemplate)
}

// handleDeleteTemplate deletes a template.
func handleDeleteTemplate(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.Template{}
	)
	id, err := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err != nil || id == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Invalid template `id`.", nil, envelope.InputError)
	}
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.errorParsing", "name", "{globals.terms.request}"), nil, envelope.InputError)
	}
	if err = app.tmpl.Delete(id); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}
