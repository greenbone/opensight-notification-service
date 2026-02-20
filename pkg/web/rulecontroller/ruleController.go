// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulecontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/rulerepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/ruleservice"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type RuleService interface {
	Get(ctx context.Context, id string) (models.Rule, error)
	List(ctx context.Context) ([]models.Rule, error)
	Create(ctx context.Context, rule models.Rule) (models.Rule, error)
	Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error)
	Delete(ctx context.Context, id string) error
}

type RuleController struct {
	ruleService RuleService
}

func NewRuleController(
	router gin.IRouter,
	ruleService RuleService,
	auth gin.HandlerFunc,
	registry *errmap.Registry,
) *RuleController {
	ctrl := &RuleController{
		ruleService: ruleService,
	}
	ctrl.RegisterRoutes(router, auth)
	ctrl.configureMappings(registry)

	return ctrl
}

func (c *RuleController) RegisterRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/rules").
		Use(middleware.AuthorizeRoles(auth, middleware.AdminRole)...)
	group.GET("/:id", c.GetRule)
	group.POST("", c.CreateRule)
	group.PUT("/:id", c.UpdateRule)
	group.DELETE("/:id", c.DeleteRule)
	group.GET("", c.ListRules)
}

func (c *RuleController) configureMappings(r *errmap.Registry) {
	r.Register(
		ruleservice.ErrRuleLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse(translation.RuleLimitReached),
	)
	r.Register(
		ruleservice.ErrRecipientRequired,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse("", "",
			map[string]string{"trigger.action.recipient": translation.RecipientRequiredForChannel}),
	)
	r.Register(
		ruleservice.ErrRecipientNotSupported,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse("", "",
			map[string]string{"trigger.action.recipient": translation.RecipientNotSupportedForChannel}),
	)
	r.Register(
		ruleservice.ErrChannelNotFound,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse("", "",
			map[string]string{"action.channel.id": translation.ChannelNotFound},
		),
	)
	r.Register(
		rulerepository.ErrInvalidID,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse(translation.InvalidID, "", nil),
	)
	r.Register(
		ruleservice.ErrOriginsNotFound,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse("", "",
			map[string]string{"trigger.origins": translation.OriginsNotFound},
		),
	)
	r.Register(
		rulerepository.ErrDuplicateRuleName,
		http.StatusBadRequest,
		errorResponses.NewErrorValidationResponse("", "",
			map[string]string{"name": translation.RuleNameAlreadyExists},
		),
	)
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
//	@Success		201		{object}	models.Rule
//	@Failure		400		{object}	errorResponses.ErrorResponse
//	@Failure		422		{object}	errorResponses.ErrorResponse
//	@Header			all		{string}	api-version	"API version"
//	@Router			/rules [post]
func (c *RuleController) CreateRule(gc *gin.Context) {
	var rule models.Rule
	if !ginEx.BindAndValidateBody(gc, &rule) {
		return
	}

	createdRule, err := c.ruleService.Create(gc.Request.Context(), rule)
	if ginEx.AddError(gc, err) {
		return
	}

	gc.JSON(http.StatusCreated, createdRule)
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
//	@Success		200		{object}	models.Rule
//	@Failure		400		{object}	errorResponses.ErrorResponse
//	@Failure		404		{object}	errorResponses.ErrorResponse
//	@Header			all		{string}	api-version	"API version"
//	@Router			/rules/{id} [put]
func (c *RuleController) UpdateRule(gc *gin.Context) {
	id := gc.Param("id")

	var rule models.Rule
	if !ginEx.BindAndValidateBody(gc, &rule) {
		return
	}

	updatedRule, err := c.ruleService.Update(gc.Request.Context(), id, rule)
	if ginEx.AddError(gc, err) {
		return
	}

	gc.JSON(http.StatusOK, updatedRule)
}

// DeleteRule
//
//	@Summary		Delete Rule
//	@Description	Delete a rule.
//	@Tags			rule
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"unique ID of the rule"
//	@Success		204	"deleted"
//	@Failure		400 {object}	errorResponses.ErrorResponse "invalid id"
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rules/{id} [delete]
func (c *RuleController) DeleteRule(gc *gin.Context) {
	id := gc.Param("id")
	err := c.ruleService.Delete(gc.Request.Context(), id)
	if ginEx.AddError(gc, err) {
		return
	}

	gc.Status(http.StatusNoContent)
}

// GetRule
//
//	@Summary		Get a rule by id
//	@Description	Returns the rule
//	@Tags			rule
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path		string	true	"unique id of the rule"
//	@Success		200	{object}	models.Rule
//	@Failure		400 {object}	errorResponses.ErrorResponse "invalid id"
//	@Failure		404	{object}	errorResponses.ErrorResponse
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rulse/{id} [get]
func (c *RuleController) GetRule(gc *gin.Context) {
	id := gc.Param("id")
	rule, err := c.ruleService.Get(gc.Request.Context(), id)
	if ginEx.AddError(gc, err) {
		return
	}

	gc.JSON(http.StatusOK, rule)
}

// ListRules
//
//	@Summary		List Rules
//	@Description	Returns a list of rules matching the provided filters.
//	@Tags			rule
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{object}	[]models.Rule
//	@Header			all	{string}	api-version	"API version"
//	@Router			/rules [get]
func (c *RuleController) ListRules(gc *gin.Context) {
	rules, err := c.ruleService.List(gc.Request.Context())
	if ginEx.AddError(gc, err) {
		return
	}
	if len(rules) == 0 { // return empty array rather than null
		rules = []models.Rule{}
	}

	gc.JSON(http.StatusOK, rules)
}
