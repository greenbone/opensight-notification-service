package mattermostcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/mattermostdto"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type MattermostController struct {
	notificationChannelServicer notificationchannelservice.NotificationChannelService
	mattermostChannelService    notificationchannelservice.MattermostChannelService
}

func NewMattermostController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.NotificationChannelService,
	mattermostChannelService notificationchannelservice.MattermostChannelService,
	auth gin.HandlerFunc,
	registry *errmap.Registry,
) *MattermostController {
	ctrl := &MattermostController{
		notificationChannelServicer: notificationChannelServicer,
		mattermostChannelService:    mattermostChannelService,
	}
	ctrl.registerRoutes(router, auth)
	ctrl.configureMappings(registry)
	return ctrl
}

func (mc *MattermostController) configureMappings(r *errmap.Registry) {
	r.Register(
		notificationchannelservice.ErrMattermostChannelLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Mattermost channel limit reached."),
	)
	r.Register(
		notificationchannelservice.ErrListMattermostChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
	r.Register(
		notificationchannelservice.ErrMattermostChannelNameExists,
		http.StatusBadRequest,
		// TODO: 27.01.2026 stolksdorf - who do I get the message from notificationchannelservice.ErrMattermostChannelNameExists
		errorResponses.NewErrorGenericResponse("Channel name should be unique."),
	)
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
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	mattermostChannel, err := mc.mattermostChannelService.CreateMattermostChannel(c, channel)
	if err != nil {
		ginEx.AddError(c, err)
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
		ginEx.AddError(c, err)
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
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	notificationChannel := mattermostdto.MapMattermostToNotificationChannel(channel)
	updated, err := mc.notificationChannelServicer.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		ginEx.AddError(c, err)
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
		ginEx.AddError(c, err)
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
//	@Success		204 "Mattermost server test message sent successfully"
//	@Failure		400			{object}	map[string]string
//	@Router			/notification-channel/mattermost/check [post]
func (mc *MattermostController) SendMattermostTestMessage(c *gin.Context) {
	var channel mattermostdto.MattermostNotificationChannelCheckRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	if err := mc.mattermostChannelService.SendMattermostTestMessage(channel.WebhookUrl); err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
