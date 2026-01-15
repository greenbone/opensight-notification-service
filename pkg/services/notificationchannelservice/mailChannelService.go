package notificationchannelservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

type MailChannelServicer interface {
	CreateMailChannel(c *gin.Context, channel models.MailNotificationChannel)
}

type MailChannelService struct {
	notificationChannelService NotificationChannelServicer
}

func NewMailChannelService(notificationChannelService NotificationChannelServicer) *MailChannelService {
	return &MailChannelService{notificationChannelService: notificationChannelService}
}

func (v *MailChannelService) ValidateFields(channel models.MailNotificationChannel) *errs.ErrorResponse {
	errors := make(errs.FieldErrors)
	if channel.Domain == nil || *channel.Domain == "" {
		errors["domain"] = "Domain cannot be empty."
	}
	if channel.Port == nil {
		errors["port"] = "Port cannot be empty."
	}
	if channel.SenderEmailAddress == nil || *channel.SenderEmailAddress == "" {
		errors["senderEmailAddress"] = "Sender email address cannot be empty."
	}
	if channel.ChannelName == nil || *channel.ChannelName == "" {
		errors["channelName"] = "Channel Name cannot be empty."
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

func (v *MailChannelService) MailChannelAlreadyExists(c *gin.Context) *errs.ErrorResponse {
	channels, err := v.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		return &errs.ErrorResponse{
			Type:   "greenbone/generic-error",
			Title:  "Internal server error",
			Errors: nil,
		}
	}

	if len(channels) > 0 {
		return &errs.ErrorResponse{
			Type:   "greenbone/generic-error",
			Title:  "Mail channel already exists.",
			Errors: map[string]string{"channelName": "Mail channel already exists."},
		}
	}
	return nil
}

func (v *MailChannelService) CreateMailChannel(c *gin.Context, channel models.MailNotificationChannel) {
	if resp := v.ValidateFields(channel); resp != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if errResp := v.MailChannelAlreadyExists(c); errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	created, err := v.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := mapper.MapNotificationChannelToMail(created)
	c.JSON(http.StatusCreated, response)
}
