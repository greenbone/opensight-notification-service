package mattermostcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
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
	ctrl := &MattermostController{
		service:                  service,
		mattermostChannelService: mattermostChannelService,
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

	err := c.ShouldBindJSON(&channel)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	mattermostChannel, err := mc.mattermostChannelService.CreateMattermostChannel(c, channel)
	if err != nil {
		_ = c.Error(err)
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
		_ = c.Error(err)
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
		_ = c.Error(err)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	notificationChannel := mapper.MapMattermostToNotificationChannel(channel)
	updated, err := mc.service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		_ = c.Error(err)
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
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
