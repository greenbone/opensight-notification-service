package mailcontroller

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

type MailController struct {
	Service            port.NotificationChannelService
	MailChannelService port.MailChannelService
}

func NewMailController(router gin.IRouter, service port.NotificationChannelService, mailChannelService port.MailChannelService, auth gin.HandlerFunc) *MailController {
	ctrl := &MailController{Service: service, MailChannelService: mailChannelService}
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
}

// CreateMailChannel
//
//	@Summary		Create Mail Channel
//	@Description	Create a new mail notification channel
//	@Tags			mail-channel
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailChannel	body		models.MailNotificationChannelRequest	true	"Mail channel to add"
//	@Success		201			{object}	models.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail [post]
func (mc *MailController) CreateMailChannel(c *gin.Context) {
	var channel request.MailNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrMailChannelBadRequest)
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
//	@Success		200		{array}		models.MailNotificationChannelRequest
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mail [get]
func (mc *MailController) ListMailChannelsByType(c *gin.Context) {
	channels, err := mc.Service.ListNotificationChannelsByType(c, models.ChannelTypeMail)

	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	c.JSON(http.StatusOK, mapper.MapNotificationChannelsToMail(channels))
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
//	@Param			MailChannel	body		models.MailNotificationChannelRequest	true	"Mail channel to update"
//	@Success		200			{object}	models.MailNotificationChannelRequest
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail/{id} [put]
func (mc *MailController) UpdateMailChannel(c *gin.Context) {
	id := c.Param("id")
	var channel request.MailNotificationChannelRequest
	if err := c.ShouldBindJSON(&channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, notificationchannelservice.ErrMailChannelBadRequest)
		return
	}

	if err := mc.validateFields(channel); err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "Mandatory fields of mail configuration cannot be empty",
			err, nil)
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	updated, err := mc.Service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		restErrorHandler.NotificationChannelErrorHandler(c, "", nil, err)
		return
	}
	response := mapper.MapNotificationChannelToMail(updated)
	c.JSON(http.StatusOK, response)
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

func (v *MailController) validateFields(channel request.MailNotificationChannelRequest) map[string]string {
	errors := make(map[string]string)
	if channel.Domain == "" {
		errors["domain"] = "A Mailhub is required."
	}
	if channel.Port == "" {
		errors["port"] = "A port is required."
	}
	if channel.SenderEmailAddress == "" {
		errors["senderEmailAddress"] = "A sender is required."
	} else {
		v.validateEmailAddress(channel.SenderEmailAddress, errors)
	}
	if channel.ChannelName == "" {
		errors["channelName"] = "A Channel Name is required."
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (v *MailController) validateEmailAddress(channel string, errors map[string]string) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, channel)
	if !matched {
		errors["senderEmailAddress"] = "A sender is required."
	}
}
