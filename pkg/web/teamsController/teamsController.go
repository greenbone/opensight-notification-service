package teamsController

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamsController/dto"
)

var ErrTeamsChannelBadRequest = errors.New("bad request for teams channel")

type TeamsController struct {
	notificationChannelServicer notificationchannelservice.NotificationChannelService
	teamsChannelService         notificationchannelservice.TeamsChannelService
}

func NewTeamsController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.NotificationChannelService,
	teamsChannelService notificationchannelservice.TeamsChannelService,
	auth gin.HandlerFunc,
) *TeamsController {
	ctrl := &TeamsController{
		notificationChannelServicer: notificationChannelServicer,
		teamsChannelService:         teamsChannelService,
	}
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
	group.POST("/check", tc.SendTeamsTestMessage)

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
	var channel dto.TeamsNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrTeamsChannelBadRequest)
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
	channels, err := tc.notificationChannelServicer.ListNotificationChannelsByType(c, models.ChannelTypeTeams)

	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	c.JSON(http.StatusOK, dto.MapNotificationChannelsToTeams(channels))
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
	var channel dto.TeamsNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrTeamsChannelBadRequest)
		return
	}

	if err := tc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of teams configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := dto.MapTeamsToNotificationChannel(channel)
	updated, err := tc.notificationChannelServicer.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	response := dto.MapNotificationChannelToTeams(updated)
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
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (tc *TeamsController) validateFields(channel dto.TeamsNotificationChannelRequest) map[string]string {
	errs := make(map[string]string)
	if channel.ChannelName == "" {
		errs["channelName"] = "A channel name is required."
	}
	if channel.WebhookUrl == "" {
		errs["webhookUrl"] = "A Webhook URL is required."
	} else {
		var re = regexp.MustCompile(`^https://[\w.-]+/webhook/[a-zA-Z0-9]+$`)
		if !re.MatchString(channel.WebhookUrl) {
			errs["webhookUrl"] = "Invalid teams webhook URL format."
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// SendTeamsTestMessage
//
//	@Summary		Check Teams server
//	@Description	Check if a Teams server is able to send a test message
//	@Tags			Teams-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"Teams channel ID"
//	@Success		204 "Teams server test message sent successfully"
//	@Failure		400			{object}	map[string]string
//	@Router			/notification-channel/Teams/check [post]
func (tc *TeamsController) SendTeamsTestMessage(c *gin.Context) {
	var channel dto.TeamsNotificationChannelCheckRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrTeamsChannelBadRequest)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	if err := tc.teamsChannelService.SendTeamsTestMessage(channel.WebhookUrl); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Failed to send test message to Teams server", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}
