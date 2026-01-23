package teamsController

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type TeamsController struct {
	service             port.NotificationChannelService
	teamsChannelService port.TeamsChannelService
}

func NewTeamsController(
	router gin.IRouter,
	service port.NotificationChannelService, teamsChannelService port.TeamsChannelService,
	auth gin.HandlerFunc,
) *TeamsController {
	ctrl := &TeamsController{service: service, teamsChannelService: teamsChannelService}
	ctrl.registerRoutes(router, auth)
	return ctrl
}

func (tc *TeamsController) registerRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/notification-channel/teams").
		Use(middleware.AuthorizeRoles(auth, "admin")...)
	group.POST("", tc.CreateTeamsChannel)
	group.GET("", tc.ListTeamsChannelsByType)
	group.PUT("/:id", tc.UpdateTeamsChannel)
	group.DELETE("/:id", tc.DeleteTeamsChannel)
}

// CreateTeamsChannel
//
//	@Summary		Create Teams Channel
//	@Description	Create a new teams notification channel
//	@Tags			teams-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			TeamsChannel	body		request.TeamsNotificationChannelRequest	true	"Teams channel to add"
//	@Success		201			{object}	request.TeamsNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/teams [post]
func (tc *TeamsController) CreateTeamsChannel(c *gin.Context) {
	var channel request.TeamsNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrTeamsChannelBadRequest)
		return
	}

	if err := tc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of teams configuration cannot be empty",
			err, nil)
		return
	}

	teamsChannel, err := tc.teamsChannelService.CreateTeamsChannel(c, channel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}

	c.JSON(http.StatusCreated, teamsChannel)
}

// ListTeamsChannelsByType
//
//	@Summary		List Teams Channels
//	@Description	List teams notification channels by type
//	@Tags			teams-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		request.TeamsNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/teams [get]
func (tc *TeamsController) ListTeamsChannelsByType(c *gin.Context) {
	channels, err := tc.service.ListNotificationChannelsByType(c, models.ChannelTypeTeams)

	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	c.JSON(http.StatusOK, mapper.MapNotificationChannelsToTeams(channels))
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
//	@Param			TeamsChannel	body		request.TeamsNotificationChannelRequest	true	"Teams channel to update"
//	@Success		200			{object}	request.TeamsNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		404 		{object}    map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/teams/{id} [put]
func (tc *TeamsController) UpdateTeamsChannel(c *gin.Context) {
	id := c.Param("id")
	var channel request.TeamsNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrTeamsChannelBadRequest)
		return
	}

	if err := tc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of teams configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := mapper.MapTeamsToNotificationChannel(channel)
	updated, err := tc.service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	response := mapper.MapNotificationChannelToTeams(updated)
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
	if err := tc.service.DeleteNotificationChannel(c, id); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (tc *TeamsController) validateFields(channel request.TeamsNotificationChannelRequest) map[string]string {
	errors := make(map[string]string)
	if channel.WebhookUrl == "" {
		errors["webhookUrl"] = "A WebhookUrl is required."
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
