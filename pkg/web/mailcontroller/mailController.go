package mailcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type MailController struct {
	Service notificationchannelservice.NotificationChannelServicer
}

func NewMailController(router gin.IRouter, service notificationchannelservice.NotificationChannelServicer, auth gin.HandlerFunc) *MailController {
	ctrl := &MailController{Service: service}
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
//	@Param			MailChannel	body		models.MailNotificationChannel	true	"Mail channel to add"
//	@Success		201			{object}	models.MailNotificationChannel
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail [post]
func (mc *MailController) CreateMailChannel(c *gin.Context) {
	var channel models.MailNotificationChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp := mc.validateMailChannelFields(channel); resp != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	created, err := mc.Service.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := mapper.MapNotificationChannelToMail(created)
	c.JSON(http.StatusCreated, response)
}

// ListMailChannelsByType
//
//	@Summary		List Mail Channels
//	@Description	List mail notification channels by type
//	@Tags			mail-channel
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			type	query		string	false	"Channel type"
//	@Success		200		{array}		models.MailNotificationChannel
//	@Failure		500		{object}	map[string]string
//	@Router			/notification-channel/mail [get]
func (mc *MailController) ListMailChannelsByType(c *gin.Context) {
	channels, err := mc.Service.ListNotificationChannelsByType(c, "mail")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
//	@Param			MailChannel	body		models.MailNotificationChannel	true	"Mail channel to update"
//	@Success		200			{object}	models.MailNotificationChannel
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/notification-channel/mail/{id} [put]
func (mc *MailController) UpdateMailChannel(c *gin.Context) {
	id := c.Param("id")
	var channel models.MailNotificationChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp := mc.validateMailChannelFields(channel); resp != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	updated, err := mc.Service.UpdateNotificationChannel(c, id, notificationChannel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (mc *MailController) validateMailChannelFields(channel models.MailNotificationChannel) *errs.ErrorResponse {
	errors := make(errs.FieldErrors)
	if channel.Domain == nil || *channel.Domain == "" {
		errors["domain"] = "Domain cannot be empty."
	}
	if channel.Port == nil {
		errors["port"] = "Port cannot be empty."
	}
	if channel.SenderEmailAddress == nil || *channel.SenderEmailAddress == "" {
		errors["sender_email_address"] = "Sender email address cannot be empty."
	}
	if channel.ChannelName == nil || *channel.ChannelName == "" {
		errors["channel_name"] = "Channel Name cannot be empty."
	}

	if len(errors) > 0 {
		return &errs.ErrorResponse{
			Type:   "greenbone/generic-error",
			Title:  "Mandatory fields of mail configuration cannot be empty",
			Errors: errors,
		}
	}
	return nil
}
