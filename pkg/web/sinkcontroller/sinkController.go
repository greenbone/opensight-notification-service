// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package sinkcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
)

type SinkController struct{}

// GetSink
//
//	@Summary		Get a sink by id
//	@Description	Returns the sink
//	@Tags			sink
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path		string	true	"unique id of the sink"
//	@Success		200	{object}	query.ResponseWithMetadata[models.Sink]
//	@Failure		404	{object}	errorResponses.ErrorResponse
//	@Header			all	{string}	api-version	"API version"
//	@Router			/sinks/{id} [get]
func (c *SinkController) GetSink(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// Create Sink
//
//	@Summary		Create Sink
//	@Description	Creates a new sink. E.g. sending a mail or a message to a mattermost channel etc.
//	@Tags			sink
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth[admin]
//	@Param			sink	body		models.Sink	true	"sink to create"
//	@Success		201		{object}	query.ResponseWithMetadata[models.Sink]
//	@Failure		409		{object}	errorResponses.ErrorResponse	"duplicate"
//	@Header			all		{string}	api-version						"API version"
//	@Router			/sinks [post]
func (c *SinkController) CreateSink(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// Update Sink
//
//	@Summary		Update Sink
//	@Description	Updates an sink.
//	@Tags			sink
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth[admin]
//	@Param			id		path		string		true	"unique ID of sink"
//	@Param			sink	body		models.Sink	true	"update/replace sink with given one"
//	@Success		201		{object}	query.ResponseWithMetadata[models.Sink]
//	@Failure		409		{object}	errorResponses.ErrorResponse	"duplicate"
//	@Header			all		{string}	api-version						"API version"
//	@Router			/sinks/{id} [put]
func (c *SinkController) UpdateSink(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// DeleteSink
//
//	@Summary		Delete Sink
//	@Description	Deletes an sink.
//	@Tags			sink
//	@Security		KeycloakAuth[admin]
//	@Param			id	path	string	true	"unique ID of sink"
//	@Success		204	"deleted"
//	@Header			all	{string}	api-version
//	@Router			/sink/{id} [delete]
func (c *SinkController) DeleteSink(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// ListSinks
//
//	@Summary		List Sinks
//	@Description	Returns a list of sinks matching the provided filters.
//	@Tags			sink
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MatchCriterias	body		query.ResultSelector	true	"filters, paging and sorting"
//	@Success		200				{object}	query.ResponseListWithMetadata[models.Sink]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/sinks [put]
func (c *SinkController) ListSinks(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// GetOptions
//
//	@Summary		ListSinks filter options
//	@Description	Get filter options for listing sinks.
//	@Tags			sink
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{object}	query.ResponseWithMetadata[[]query.FilterOption]
//	@Header			all	{string}	api-version	"API version"
//	@Router			/sinks/options [get]
func (c *SinkController) GetOptions(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}
