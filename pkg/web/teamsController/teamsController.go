package teamsController

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamsController/teamsdto"
)

var ErrTeamsChannelBadRequest = errors.New("bad request for teams channel")

type TeamsController struct {
	notificationChannelServicer notificationchannelservice.NotificationChannelService
	teamsChannelService         notificationchannelservice.TeamsChannelService
}

func AddTeamsController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.NotificationChannelService,
	teamsChannelService notificationchannelservice.TeamsChannelService,
	auth gin.HandlerFunc,
	registry *errmap.Registry,

) *TeamsController {
	ctrl := &TeamsController{
		notificationChannelServicer: notificationChannelServicer,
		teamsChannelService:         teamsChannelService,
	}

	group := router.Group("/notification-channel/teams").
		Use(middleware.AuthorizeRoles(auth, "admin")...)

	group.POST("", ctrl.CreateTeamsChannel)
	group.GET("", ctrl.ListTeamsChannels)
	group.PUT("/:id", ctrl.UpdateTeamsChannel)
	group.DELETE("/:id", ctrl.DeleteTeamsChannel)
	group.POST("/check", ctrl.SendTeamsTestMessage)

	router.Use(errorHandler(gin.ErrorTypePrivate))
	ctrl.configureMappings(registry)
	return ctrl
}

func errorHandler(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, errorValue := range c.Errors.ByType(errorType) {
			if errors.Is(errorValue, notificationchannelservice.ErrTeamsMassageDelivery) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorGenericResponse(errorValue.Error()))
				return
			}
		}
	}
}

func (tc *TeamsController) configureMappings(r *errmap.Registry) {
	r.Register(
		notificationchannelservice.ErrTeamsChannelLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse(notificationchannelservice.ErrTeamsChannelLimitReached.Error()),
	)
	r.Register(
		notificationchannelservice.ErrListTeamsChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
	r.Register(
		notificationchannelservice.ErrTeamsChannelNameExists,
		http.StatusBadRequest,
		errorResponses.NewErrorGenericResponse(notificationchannelservice.ErrTeamsChannelNameExists.Error()),
	)
}

// CreateTeamsChannel
//
//	@Summary		Create Teams Channel
//	@Description	Create a new teams notification channel
//	@Tags			teams-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			TeamsChannel	body		teamsdto.TeamsNotificationChannelRequest	true	"Teams channel to add"
//	@Success		201			{object}	teamsdto.TeamsNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/teams [post]
func (tc *TeamsController) CreateTeamsChannel(c *gin.Context) {
	var channel teamsdto.TeamsNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	teamsChannel, err := tc.teamsChannelService.CreateTeamsChannel(c, channel)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.JSON(http.StatusCreated, teamsChannel)
}

// ListTeamsChannels
//
//	@Summary		List Teams Channels
//	@Description	List teams notification channels by type
//	@Tags			teams-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		teamsdto.TeamsNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/teams [get]
func (tc *TeamsController) ListTeamsChannels(c *gin.Context) {
	channels, err := tc.notificationChannelServicer.ListNotificationChannelsByType(c, models.ChannelTypeTeams)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.JSON(http.StatusOK, teamsdto.MapNotificationChannelsToTeams(channels))
}

// UpdateTeamsChannel
//
//	@Summary		Update Teams Channel
//	@Description	Update an existing teams notification channel
//	@Tags			teams-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id			path		string						true	"Teams channel ID"
//	@Param			TeamsChannel	body		teamsdto.TeamsNotificationChannelRequest	true	"Teams channel to update"
//	@Success		200			{object}	teamsdto.TeamsNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		404 		{object}    map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/teams/{id} [put]
func (tc *TeamsController) UpdateTeamsChannel(c *gin.Context) {
	id := c.Param("id")
	var channel teamsdto.TeamsNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	notificationChannel := teamsdto.MapTeamsToNotificationChannel(channel)
	updated, err := tc.notificationChannelServicer.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}
	response := teamsdto.MapNotificationChannelToTeams(updated)
	c.JSON(http.StatusOK, response)
}

// DeleteTeamsChannel
//
//		@Summary		Delete Teams Channel
//		@Description	Delete a teams notification channel
//		@Tags			teams-channel
//		@Security		KeycloakAuth
//		@Param			id	path	string	true	"Teams channel ID"
//		@Success		204	"Deleted successfully"
//		@Failure		500	{object}	map[string]string
//	    @Failure		404 {object}    map[string]string
//		@Router			/notification-channel/teams/{id} [delete]
func (tc *TeamsController) DeleteTeamsChannel(c *gin.Context) {
	id := c.Param("id")
	if err := tc.notificationChannelServicer.DeleteNotificationChannel(c, id); err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// SendTeamsTestMessage
//
//	@Summary		Check Teams server
//	@Description	Check if a Teams server is able to send a test message
//	@Tags			teams-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"Teams channel ID"
//	@Success		204 "Teams server test message sent successfully"
//	@Failure		400			{object}	map[string]string
//	@Router			/notification-channel/teams/check [post]
func (tc *TeamsController) SendTeamsTestMessage(c *gin.Context) {
	var channel teamsdto.TeamsNotificationChannelCheckRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	if err := tc.teamsChannelService.SendTeamsTestMessage(channel.WebhookUrl); err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
