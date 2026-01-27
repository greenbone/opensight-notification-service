package mailcontroller

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/maildto"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

var ErrMailChannelBadRequest = errors.New("bad request for mail channel")

type MailController struct {
	Service            notificationchannelservice.NotificationChannelService
	MailChannelService notificationchannelservice.MailChannelService
}

func NewMailController(
	router gin.IRouter,
	service notificationchannelservice.NotificationChannelService,
	mailChannelService notificationchannelservice.MailChannelService,
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
//	@Param			MailChannel	body		maildto.MailNotificationChannelRequest	true	"Mail channel to add"
//	@Success		201			{object}	maildto.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail [post]
func (mc *MailController) CreateMailChannel(c *gin.Context) {
	var channel maildto.MailNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMailChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mail configuration cannot be empty",
			err, nil)
		return
	}

	mailChannel, err := mc.MailChannelService.CreateMailChannel(c, channel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
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
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
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
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMailChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mail configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := maildto.MapMailToNotificationChannel(channel)
	updated, err := mc.Service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
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
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
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
	if err := c.ShouldBindJSON(&mailServer); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, ErrMailChannelBadRequest)
		return
	}
	if err := mailServer.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	err := mc.MailChannelService.CheckNotificationChannelEntityConnectivity(context.Background(), id, mailServer.ToModel())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"type":  "greenbone/generic-error",
			"title": err.Error()},
		)
		return
	}

	c.Status(http.StatusNoContent)
}

func (v *MailController) validateFields(channel maildto.MailNotificationChannelRequest) map[string]string {
	errs := make(map[string]string)
	if channel.Domain == "" {
		errs["domain"] = "A Mailhub is required."
	}
	if channel.Port == 0 {
		errs["port"] = "A port is required."
	}

	if channel.IsAuthenticationRequired {
		if channel.Username != nil && *channel.Username == "" {
			errs["username"] = "Username is required."
		}
		if channel.Password != nil && *channel.Password == "" {
			errs["password"] = "Password is required."
		}
	}

	v.validateEmailAddress(channel.SenderEmailAddress, errs)
	if channel.ChannelName == "" {
		errs["channelName"] = "A Channel Name is required."
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *MailController) validateEmailAddress(senderEmailAddress string, errors map[string]string) {
	if senderEmailAddress == "" {
		errors["senderEmailAddress"] = "A sender is required."
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, senderEmailAddress)
	if !matched {
		errors["senderEmailAddress"] = "A valid sender email is required."
	}
}
