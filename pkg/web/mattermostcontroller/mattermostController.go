package mattermostcontroller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
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

	group := router.Group("/notification-channel/mattermost").
		Use(middleware.AuthorizeRoles(auth, iam.Admin)...)
	group.Use(errorHandler(gin.ErrorTypePrivate))

	group.POST("", ctrl.createMattermostChannel)
	group.GET("", ctrl.listMattermostChannels)
	group.PUT("/:id", ctrl.updateMattermostChannel)
	group.DELETE("/:id", ctrl.deleteMattermostChannel)
	group.POST("/check", ctrl.sendMattermostTestMessage)

	ctrl.configureMappings(registry)
	return ctrl
}

func errorHandler(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, errorValue := range c.Errors.ByType(errorType) {
			if errors.Is(errorValue, notificationchannelservice.ErrMattermostMassageDelivery) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorGenericResponse(errorValue.Error()))
				return
			}
		}
	}
}

func (mc *MattermostController) configureMappings(r *errmap.Registry) {
	r.Register(
		notificationchannelservice.ErrMattermostChannelLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse(translation.MattermostChannelLimitReached),
	)
	r.Register(
		notificationchannelservice.ErrListMattermostChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
	r.Register(
		notificationchannelservice.ErrMattermostChannelNameExists,
		http.StatusBadRequest,
		errorResponses.NewErrorGenericResponse(translation.MattermostChannelNameAlreadyExist),
	)
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
func (mc *MattermostController) createMattermostChannel(c *gin.Context) {
	var channel mattermostdto.MattermostNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	mattermostChannel, err := mc.mattermostChannelService.CreateMattermostChannel(c, channel)
	if ginEx.AddError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, mattermostChannel)
}

// ListMattermostChannels
//
//	@Summary		List Mattermost Channels
//	@Description	List mattermost notification channels
//	@Tags			mattermost-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200		{array}		mattermostdto.MattermostNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mattermost [get]
func (mc *MattermostController) listMattermostChannels(c *gin.Context) {
	channels, err := mc.notificationChannelServicer.ListNotificationChannelsByType(c, models.ChannelTypeMattermost)
	if ginEx.AddError(c, err) {
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
func (mc *MattermostController) updateMattermostChannel(c *gin.Context) {
	id := c.Param("id")

	var channel mattermostdto.MattermostNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	updated, err := mc.mattermostChannelService.UpdateMattermostChannel(c, id, channel)
	if ginEx.AddError(c, err) {
		return
	}

	c.JSON(http.StatusOK, updated)
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
func (mc *MattermostController) deleteMattermostChannel(c *gin.Context) {
	id := c.Param("id")

	err := mc.notificationChannelServicer.DeleteNotificationChannel(c, id)
	if ginEx.AddError(c, err) {
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
func (mc *MattermostController) sendMattermostTestMessage(c *gin.Context) {
	var channel mattermostdto.MattermostNotificationChannelCheckRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	err := mc.mattermostChannelService.SendMattermostTestMessage(channel.WebhookUrl)
	if ginEx.AddError(c, err) {
		return
	}

	c.Status(http.StatusNoContent)
}
