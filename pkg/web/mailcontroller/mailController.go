package mailcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/dtos"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type MailController struct {
	Service            port.NotificationChannelService
	MailChannelService port.MailChannelService
}

func NewMailController(
	router gin.IRouter,
	service port.NotificationChannelService, mailChannelService port.MailChannelService,
	auth gin.HandlerFunc,
) *MailController {
	ctrl := &MailController{
		Service:            service,
		MailChannelService: mailChannelService,
	}
	ctrl.registerRoutes(router, auth)
	return ctrl
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
//	@Param			MailChannel	body		request.MailNotificationChannelRequest	true	"Mail channel to add"
//	@Success		201			{object}	request.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail [post]
func (mc *MailController) CreateMailChannel(c *gin.Context) {
	var channel request.MailNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		_ = c.Error(err)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	mailChannel, err := mc.MailChannelService.CreateMailChannel(c, channel)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, mailChannel.WithEmptyPassword())
}

// ListMailChannelsByType
//
//	@Summary		List Mail Channels
//	@Description	List mail notification channels by type
//	@Tags			mail-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		request.MailNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mail [get]
func (mc *MailController) ListMailChannelsByType(c *gin.Context) {
	channels, err := mc.Service.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mapper.MapNotificationChannelsToMailWithEmptyPassword(channels))
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
//	@Param			MailChannel	body		request.MailNotificationChannelRequest	true	"Mail channel to update"
//	@Success		200			{object}	request.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail/{id} [put]
func (mc *MailController) UpdateMailChannel(c *gin.Context) {
	id := c.Param("id")
	var channel request.MailNotificationChannelRequest

	if err := c.ShouldBindJSON(&channel); err != nil {
		_ = c.Error(err)
		return
	}

	if err := channel.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	updated, err := mc.Service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		_ = c.Error(err)
		return
	}

	mailChannel := mapper.MapNotificationChannelToMail(updated)
	c.JSON(http.StatusOK, mailChannel.WithEmptyPassword())
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
		_ = c.Error(err)
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
//	@Param			MailServerConfig	body		dtos.CheckMailServerEntityRequest	true	"Mail server to check"
//	@Success		204 "Mail server reachable"
//	@Failure		400			{object}	map[string]string
//	@Failure		422 "Mail server error"
//	@Router			/notifications/mail/{id}/check [post]
func (mc *MailController) CheckMailServer(c *gin.Context) {
	id := c.Param("id")

	var mailServer dtos.CheckMailServerEntityRequest
	if err := c.ShouldBindJSON(&mailServer); err != nil {
		_ = c.Error(err)
		return
	}
	if err := mailServer.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	err := mc.Service.CheckNotificationChannelEntityConnectivity(context.Background(), id, mailServer.ToModel())
	if err != nil {
		_ = c.Error(err)
		return
		/*
			c.JSON(http.StatusUnprocessableEntity, gin.H{
					"type":  "greenbone/generic-error",
					"title": err.Error()},
				)
		*/
	}

	c.Status(http.StatusNoContent)
}
