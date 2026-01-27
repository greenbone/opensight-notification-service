package mailcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/maildto"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type MailController struct {
	Service            notificationchannelservice.NotificationChannelService
	MailChannelService notificationchannelservice.MailChannelService
}

func NewMailController(
	router gin.IRouter,
	service notificationchannelservice.NotificationChannelService,
	mailChannelService notificationchannelservice.MailChannelService,
	auth gin.HandlerFunc,
	registry errmap.ErrorRegistry,
) *MailController {
	ctrl := &MailController{
		Service:            service,
		MailChannelService: mailChannelService,
	}
	ctrl.registerRoutes(router, auth)
	ctrl.configureMappings(registry)
	return ctrl
}

func (mc *MailController) configureMappings(r errmap.ErrorRegistry) {
	r.Register(
		notificationchannelservice.ErrMailChannelLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Mail channel limit reached."),
	)

	r.Register(
		notificationchannelservice.ErrListMailChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
}

func (mc *MailController) registerRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/notification-channel/mail").
		Use(middleware.AuthorizeRoles(auth, "admin")...)
	group.POST("", mc.CreateMailChannel)
	group.GET("", mc.ListMailChannelsByType)
	group.PUT("/:id", mc.UpdateMailChannel)
	group.DELETE("/:id", mc.DeleteMailChannel)
	group.POST("/:id/check", mc.CheckMailServer)
}

// CreateMailChannel
//
//	@Summary		Create Mail Channel
//	@Description	Create a new mail notification channel
//	@Tags			mail-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailChannel	body		maildto.MailNotificationChannelRequest	true	"Mail channel to add"
//	@Success		201			{object}	maildto.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail [post]
func (mc *MailController) CreateMailChannel(c *gin.Context) {
	var channel maildto.MailNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	mailChannel, err := mc.MailChannelService.CreateMailChannel(c, channel)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mailChannel)
}

// ListMailChannelsByType
//
//	@Summary		List Mail Channels
//	@Description	List mail notification channels by type
//	@Tags			mail-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		maildto.MailNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mail [get]
func (mc *MailController) ListMailChannelsByType(c *gin.Context) {
	channels, err := mc.Service.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.JSON(http.StatusOK, maildto.MapNotificationChannelsToMail(channels))
}

// UpdateMailChannel
//
//	@Summary		Update Mail Channel
//	@Description	Update an existing mail notification channel
//	@Tags			mail-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			id			path		string						true	"Mail channel ID"
//	@Param			MailChannel	body		maildto.MailNotificationChannelRequest	true	"Mail channel to update"
//	@Success		200			{object}	maildto.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail/{id} [put]
func (mc *MailController) UpdateMailChannel(c *gin.Context) {
	id := c.Param("id")
	var channel maildto.MailNotificationChannelRequest
	if !ginEx.BindAndValidateBody(c, &channel) {
		return
	}

	notificationChannel := maildto.MapMailToNotificationChannel(channel)
	updated, err := mc.Service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	mailChannel := maildto.MapNotificationChannelToMail(updated)
	c.JSON(http.StatusOK, mailChannel)
}

// DeleteMailChannel
//
//	@Summary		Delete Mail Channel
//	@Description	Delete a mail notification channel
//	@Tags			mail-channel
//	@Security		KeycloakAuth
//	@Param			id	path	string	true	"Mail channel ID"
//	@Success		204	"Deleted successfully"
//	@Failure		500	{object}	map[string]string
//	@Router			/notification-channel/mail/{id} [delete]
func (mc *MailController) DeleteMailChannel(c *gin.Context) {
	id := c.Param("id")

	if err := mc.Service.DeleteNotificationChannel(c, id); err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// CheckMailServer
//
//	@Summary		Check mail server
//	@Description	Check if a mail server is reachable
//	@Tags			mailserver
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailServerConfig	body		maildto.CheckMailServerEntityRequest	true	"Mail server to check"
//	@Success		204 "Mail server reachable"
//	@Failure		400			{object}	map[string]string
//	@Failure		422 "Mail server error"
//	@Router			/notifications/mail/{id}/check [post]
func (mc *MailController) CheckMailServer(c *gin.Context) {
	id := c.Param("id")

	var mailServer maildto.CheckMailServerEntityRequest
	if !ginEx.BindAndValidateBody(c, &mailServer) {
		return
	}

	err := mc.MailChannelService.CheckNotificationChannelEntityConnectivity(context.Background(), id, mailServer.ToModel())
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
