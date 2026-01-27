package mattermostcontroller

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/mattermostdto"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

var ErrMattermostChannelBadRequest = errors.New("bad request for mattermost channel")

type MattermostController struct {
	notificationChannelServicer notificationchannelservice.NotificationChannelService
	mattermostChannelService    notificationchannelservice.MattermostChannelService
}

func NewMattermostController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.NotificationChannelService,
	mattermostChannelService notificationchannelservice.MattermostChannelService,
	auth gin.HandlerFunc,
) *MattermostController {
	ctrl := &MattermostController{
		notificationChannelServicer: notificationChannelServicer,
		mattermostChannelService:    mattermostChannelService,
	}
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
	group.POST("/check", mc.SendMattermostTestMessage)
}

// CreateMattermostChannel
//
//	@Summary		Create Mattermost Channel
//	@Description	Create a new mattermost notification channel
//	@Tags			mattermost-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MattermostChannel	body		mattermostdto.MattermostNotificationChannelRequest	true	"Mattermost channel to add"
//	@Success		201			{object}	mattermostdto.MattermostNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mattermost [post]
func (mc *MattermostController) CreateMattermostChannel(c *gin.Context) {
	var channel mattermostdto.MattermostNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMattermostChannelBadRequest)
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
//	@Success		200		{array}		mattermostdto.MattermostNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mattermost [get]
func (mc *MattermostController) ListMattermostChannelsByType(c *gin.Context) {
	channels, err := mc.notificationChannelServicer.ListNotificationChannelsByType(c, models.ChannelTypeMattermost)

	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	c.JSON(http.StatusOK, mattermostdto.MapNotificationChannelsToMattermost(channels))
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
//	@Param			MattermostChannel	body		mattermostdto.MattermostNotificationChannelRequest	true	"Mattermost channel to update"
//	@Success		200			{object}	mattermostdto.MattermostNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		404 		{object}    map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mattermost/{id} [put]
func (mc *MattermostController) UpdateMattermostChannel(c *gin.Context) {
	id := c.Param("id")
	var channel mattermostdto.MattermostNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMattermostChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mattermost configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := mattermostdto.MapMattermostToNotificationChannel(channel)
	updated, err := mc.notificationChannelServicer.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	response := mattermostdto.MapNotificationChannelToMattermost(updated)
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
	if err := mc.notificationChannelServicer.DeleteNotificationChannel(c, id); err != nil {
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
//	@Router			/notification-channel/mattermost/check [post]
func (mc *MattermostController) SendMattermostTestMessage(c *gin.Context) {
	var channel mattermostdto.MattermostNotificationChannelCheckRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMattermostChannelBadRequest)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	if err := mc.mattermostChannelService.SendMattermostTestMessage(channel.WebhookUrl); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Failed to send test message to mattermost server", nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (v *MattermostController) validateFields(channel mattermostdto.MattermostNotificationChannelRequest) map[string]string {
	errs := make(map[string]string)
	if channel.ChannelName == "" {
		errs["channelName"] = "A channel name is required."
	}
	if channel.WebhookUrl == "" {
		errs["webhookUrl"] = "A webhook URL is required."
	} else {
		var re = regexp.MustCompile(`^https://[\w.-]+/hooks/[a-zA-Z0-9]+$`)
		if !re.MatchString(channel.WebhookUrl) {
			errs["webhookUrl"] = "Invalid mattermost webhook URL format."
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
