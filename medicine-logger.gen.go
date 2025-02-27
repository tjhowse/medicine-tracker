// Package main provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// LogID defines model for LogID.
type LogID = int

// MedicineID defines model for MedicineID.
type MedicineID = int

// MedicineLogEntry defines model for MedicineLogEntry.
type MedicineLogEntry struct {
	Count      float32    `json:"count"`
	LogId      LogID      `json:"log_id"`
	MedicineId MedicineID `json:"medicine_id"`
	Note       string     `json:"note"`
	Time       time.Time  `json:"time"`
}

// MedicineType defines model for MedicineType.
type MedicineType struct {
	Dose       float32    `json:"dose"`
	MedicineId MedicineID `json:"medicine_id"`
	Name       string     `json:"name"`
}

// UserSettings defines model for UserSettings.
type UserSettings struct {
	Name string `json:"name"`
}

// DeleteApiV1MedicineLogParams defines parameters for DeleteApiV1MedicineLog.
type DeleteApiV1MedicineLogParams struct {
	// LogId Delete entry with this log ID
	LogId LogID `form:"log_id" json:"log_id"`
}

// GetApiV1MedicineLogParams defines parameters for GetApiV1MedicineLog.
type GetApiV1MedicineLogParams struct {
	// Start Get entries after this point in time
	Start *time.Time `form:"start,omitempty" json:"start,omitempty"`

	// End Get entries before this point in time
	End *time.Time `form:"end,omitempty" json:"end,omitempty"`
}

// PostApiV1MedicineLogJSONBody defines parameters for PostApiV1MedicineLog.
type PostApiV1MedicineLogJSONBody = []MedicineLogEntry

// PostApiV1MedicinesJSONBody defines parameters for PostApiV1Medicines.
type PostApiV1MedicinesJSONBody = []MedicineType

// PostApiV1MedicineLogJSONRequestBody defines body for PostApiV1MedicineLog for application/json ContentType.
type PostApiV1MedicineLogJSONRequestBody = PostApiV1MedicineLogJSONBody

// PostApiV1MedicinesJSONRequestBody defines body for PostApiV1Medicines for application/json ContentType.
type PostApiV1MedicinesJSONRequestBody = PostApiV1MedicinesJSONBody

// PostApiV1SettingsJSONRequestBody defines body for PostApiV1Settings for application/json ContentType.
type PostApiV1SettingsJSONRequestBody = UserSettings

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/delete-user)
	GetApiV1DeleteUser(ctx echo.Context) error

	// (GET /api/v1/logout)
	GetApiV1Logout(ctx echo.Context) error

	// (DELETE /api/v1/medicine-log)
	DeleteApiV1MedicineLog(ctx echo.Context, params DeleteApiV1MedicineLogParams) error

	// (GET /api/v1/medicine-log)
	GetApiV1MedicineLog(ctx echo.Context, params GetApiV1MedicineLogParams) error

	// (POST /api/v1/medicine-log)
	PostApiV1MedicineLog(ctx echo.Context) error

	// (GET /api/v1/medicines)
	GetApiV1Medicines(ctx echo.Context) error

	// (POST /api/v1/medicines)
	PostApiV1Medicines(ctx echo.Context) error

	// (GET /api/v1/settings)
	GetApiV1Settings(ctx echo.Context) error

	// (POST /api/v1/settings)
	PostApiV1Settings(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetApiV1DeleteUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1DeleteUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1DeleteUser(ctx)
	return err
}

// GetApiV1Logout converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Logout(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Logout(ctx)
	return err
}

// DeleteApiV1MedicineLog converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1MedicineLog(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteApiV1MedicineLogParams
	// ------------- Required query parameter "log_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "log_id", ctx.QueryParams(), &params.LogId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter log_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1MedicineLog(ctx, params)
	return err
}

// GetApiV1MedicineLog converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1MedicineLog(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetApiV1MedicineLogParams
	// ------------- Optional query parameter "start" -------------

	err = runtime.BindQueryParameter("form", true, false, "start", ctx.QueryParams(), &params.Start)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter start: %s", err))
	}

	// ------------- Optional query parameter "end" -------------

	err = runtime.BindQueryParameter("form", true, false, "end", ctx.QueryParams(), &params.End)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter end: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1MedicineLog(ctx, params)
	return err
}

// PostApiV1MedicineLog converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1MedicineLog(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1MedicineLog(ctx)
	return err
}

// GetApiV1Medicines converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Medicines(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Medicines(ctx)
	return err
}

// PostApiV1Medicines converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Medicines(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Medicines(ctx)
	return err
}

// GetApiV1Settings converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Settings(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Settings(ctx)
	return err
}

// PostApiV1Settings converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Settings(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Settings(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v1/delete-user", wrapper.GetApiV1DeleteUser)
	router.GET(baseURL+"/api/v1/logout", wrapper.GetApiV1Logout)
	router.DELETE(baseURL+"/api/v1/medicine-log", wrapper.DeleteApiV1MedicineLog)
	router.GET(baseURL+"/api/v1/medicine-log", wrapper.GetApiV1MedicineLog)
	router.POST(baseURL+"/api/v1/medicine-log", wrapper.PostApiV1MedicineLog)
	router.GET(baseURL+"/api/v1/medicines", wrapper.GetApiV1Medicines)
	router.POST(baseURL+"/api/v1/medicines", wrapper.PostApiV1Medicines)
	router.GET(baseURL+"/api/v1/settings", wrapper.GetApiV1Settings)
	router.POST(baseURL+"/api/v1/settings", wrapper.PostApiV1Settings)

}
