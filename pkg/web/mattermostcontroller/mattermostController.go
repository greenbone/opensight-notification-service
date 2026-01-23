package mattermostcontroller

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type MattermostController struct {
	service                  port.NotificationChannelService
	mattermostChannelService port.MattermostChannelService
}

func NewMattermostController(
	router gin.IRouter,
	service port.NotificationChannelService, mattermostChannelService port.MattermostChannelService,
	auth gin.HandlerFunc,
) *MattermostController {
	ctrl := &MattermostController{service: service, mattermostChannelService: mattermostChannelService}
	ctrl.registerRoutes(router, auth)
	return ctrl
}

func (mc *MattermostController) registerRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/notification-channel/mattermost").
		Use(middleware.AuthorizeRoles(auth, "admin")...)
	group.POST("", mc.CreateMattermostChannel)
	group.GET("", mc.ListMattermostChannelsByType)
	group.PUT("/:id", mc.UpdateMattermostChannel)
	group.DELETE("/:id", mc.DeleteMattermostChannel)
	group.POST("/:id/check", mc.SendMattermostTestMessage)
}

// CreateMattermostChannel
//
//	@Summary		Create Mattermost Channel
//	@Description	Create a new mattermost notification channel
//	@Tags			mattermost-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MattermostChannel	body		request.MattermostNotificationChannelRequest	true	"Mattermost channel to add"
//	@Success		201			{object}	request.MattermostNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mattermost [post]
func (mc *MattermostController) CreateMattermostChannel(c *gin.Context) {
	var channel request.MattermostNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrMattermostChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mattermost configuration cannot be empty",
			err, nil)
		return
	}

	mattermostChannel, err := mc.mattermostChannelService.CreateMattermostChannel(c, channel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}

	c.JSON(http.StatusCreated, mattermostChannel)
}

// ListMattermostChannelsByType
//
//	@Summary		List Mattermost Channels
//	@Description	List mattermost notification channels by type
//	@Tags			mattermost-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		request.MattermostNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mattermost [get]
func (mc *MattermostController) ListMattermostChannelsByType(c *gin.Context) {
	channels, err := mc.service.ListNotificationChannelsByType(c, models.ChannelTypeMattermost)

	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	c.JSON(http.StatusOK, mapper.MapNotificationChannelsToMattermost(channels))
}

// UpdateMattermostChannel
//
//	@Summary		Update Mattermost Channel
//	@Description	Update an existing mattermost notification channel
//	@Tags			mattermost-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id			path		string						true	"Mattermost channel ID"
//	@Param			MattermostChannel	body		request.MattermostNotificationChannelRequest	true	"Mattermost channel to update"
//	@Success		200			{object}	request.MattermostNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		404 		{object}    map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mattermost/{id} [put]
func (mc *MattermostController) UpdateMattermostChannel(c *gin.Context) {
	id := c.Param("id")
	var channel request.MattermostNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrMattermostChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mattermost configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := mapper.MapMattermostToNotificationChannel(channel)
	updated, err := mc.service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	response := mapper.MapNotificationChannelToMattermost(updated)
	c.JSON(http.StatusOK, response)
}

// DeleteMattermostChannel
//
//		@Summary		Delete Mattermost Channel
//		@Description	Delete a mattermost notification channel
//		@Tags			mattermost-channel
//		@Security		KeycloakAuth
//		@Param			id	path	string	true	"Mattermost channel ID"
//		@Success		204	"Deleted successfully"
//		@Failure		500	{object}	map[string]string
//	    @Failure		404 {object}    map[string]string
//		@Router			/notification-channel/mattermost/{id} [delete]
func (mc *MattermostController) DeleteMattermostChannel(c *gin.Context) {
	id := c.Param("id")
	if err := mc.service.DeleteNotificationChannel(c, id); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// SendMattermostTestMessage
//
//	@Summary		Check mattermost server
//	@Description	Check if a mattermost server is able to send a test message
//	@Tags			mattermost-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"Mattermost channel ID"
//	@Success		204 "Mattermost server test message sent successfully"
//	@Failure		400			{object}	map[string]string
//	@Router			/notification-channel/mattermost/{id}/check [post]
func (mc *MattermostController) SendMattermostTestMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mattermost channel ID is required", nil,
			notificationchannelservice.ErrMattermostChannelBadRequest)
		return
	}

	if err := mc.mattermostChannelService.SendMattermostTestMessage(c, id); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Failed to send test message to mattermost server", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (v *MattermostController) validateFields(channel request.MattermostNotificationChannelRequest) map[string]string {
	errors := make(map[string]string)
	if channel.ChannelName == "" {
		errors["channelName"] = "A channel name is required."
	}
	if channel.WebhookUrl == "" {
		errors["webhookUrl"] = "A webhook URL is required."
	} else {
		var re = regexp.MustCompile(`^https://[\w.-]+/hooks/[a-zA-Z0-9]+$`)
		if !re.MatchString(channel.WebhookUrl) {
			errors["webhookUrl"] = "Invalid mattermost webhook URL format."
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
