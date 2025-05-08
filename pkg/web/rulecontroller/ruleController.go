// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulecontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
)

type RuleController struct{}

// GetRule
//
//	@Summary		Get a rule by id
//	@Description	Returns the rule
//	@Tags			rule
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path		string	true	"unique id of the rule"
//	@Success		200	{object}	query.ResponseWithMetadata[models.Rule]
//	@Failure		404	{object}	errorResponses.ErrorResponse
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rulse/{id} [get]
func (c *RuleController) GetRule(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// CreateRule
//
//	@Summary		Create Rule
//	@Description	Create a new rule. A rule determines on which conditions which action is triggered.
//	@Tags			rule
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			rule	body		models.Rule	true	"new rule"
//	@Success		201		{object}	query.ResponseWithMetadata[models.Rule]
//	@Header			all		{string}	api-version	"API version"
//	@Router			/rules [post]
func (c *RuleController) CreateRule(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// GetRuleOptions
//
//	@Summary		List available settings for rules.
//	@Description	This gives information about the possible rules. Returns a list of all available event levels, event origins for the trigger condition as well as a list of possbible sinke.
//	@Tags			rule
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{object}	query.ResponseWithMetadata[models.RuleOptions]
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rules/ruleoptions [get]
func (c *RuleController) GetRuleOptions(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// UpdateRule
//
//	@Summary		Update Rule
//	@Description	Update/replace a rule.
//	@Tags			rule
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id		path		string		true	"unique ID of the rule"
//	@Param			rule	body		models.Rule	true	"updated rule"
//	@Success		200		{object}	query.ResponseWithMetadata[models.Rule]
//	@Header			all		{string}	api-version	"API version"
//	@Router			/rules/{id} [put]
func (c *RuleController) UpdateRule(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// DeleteRule
//
//	@Summary		Delete Rule
//	@Description	Delete a rule.
//	@Tags			rule
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"unique ID of the rule"
//	@Success		204	"deleted"
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rules/{id} [delete]
func (c *RuleController) DeleteRule(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// ListRules
//
//	@Summary		List Rules
//	@Description	Returns a list of rules matching the provided filters.
//	@Tags			rule
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MatchCriterias	body		query.ResultSelector	true	"filters, paging and sorting"
//	@Success		200				{object}	query.ResponseListWithMetadata[models.Rule]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/rules [put]
func (c *RuleController) ListRules(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// GetListOptions
//
//	@Summary		ListRules filter options
//	@Description	Get filter options for listing rules.
//	@Tags			rule
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{object}	query.ResponseWithMetadata[[]query.FilterOption]
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rules/options [get]
func (c *RuleController) GetListOptions(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}
