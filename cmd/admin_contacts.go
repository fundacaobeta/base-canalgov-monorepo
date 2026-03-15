package main

import (
	"encoding/json"
	"strconv"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

type contactSegmentReq struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Filters     json.RawMessage `json:"filters"`
}

func handleGetContactSegments(r *fastglue.Request) error {
	app := r.Context.(*App)
	segments, err := app.user.GetContactSegments()
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(segments)
}

func handleGetContactSegment(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if id <= 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.invalid", "name", "`id`"), nil, envelope.InputError)
	}
	segment, err := app.user.GetContactSegment(id)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(segment)
}

func handleCreateContactSegment(r *fastglue.Request) error {
	app := r.Context.(*App)
	var req contactSegmentReq
	if err := r.Decode(&req, "json"); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.InputError, app.i18n.Ts("globals.messages.errorParsing", "name", "{globals.terms.request}"), nil))
	}

	if req.Name == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.empty", "name", "name"), nil, envelope.InputError)
	}

	segment, err := app.user.CreateContactSegment(req.Name, req.Description, req.Filters)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(segment)
}

func handleUpdateContactSegment(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if id <= 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.invalid", "name", "`id`"), nil, envelope.InputError)
	}

	var req contactSegmentReq
	if err := r.Decode(&req, "json"); err != nil {
		return sendErrorEnvelope(r, envelope.NewError(envelope.InputError, app.i18n.Ts("globals.messages.errorParsing", "name", "{globals.terms.request}"), nil))
	}

	segment, err := app.user.UpdateContactSegment(id, req.Name, req.Description, req.Filters)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(segment)
}

func handleDeleteContactSegment(r *fastglue.Request) error {
	app := r.Context.(*App)
	id, _ := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if id <= 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, app.i18n.Ts("globals.messages.invalid", "name", "`id`"), nil, envelope.InputError)
	}

	if err := app.user.DeleteContactSegment(id); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}
