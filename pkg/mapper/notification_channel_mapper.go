package mapper

import (
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/request"
)

// MapNotificationChannelToMail maps NotificationChannel to MailNotificationChannelRequest.
func MapNotificationChannelToMail(channel models.NotificationChannel) request.MailNotificationChannelRequest {
	return request.MailNotificationChannelRequest{
		Id:                       channel.Id,
		ChannelName:              *channel.ChannelName,
		Domain:                   *channel.Domain,
		Port:                     *channel.Port,
		IsAuthenticationRequired: *channel.IsAuthenticationRequired,
		IsTlsEnforced:            *channel.IsTlsEnforced,
		Username:                 channel.Username,
		Password:                 channel.Password,
		MaxEmailAttachmentSizeMb: channel.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    channel.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       *channel.SenderEmailAddress,
	}
}

func MapMailToNotificationChannel(mail request.MailNotificationChannelRequest) models.NotificationChannel {
	return models.NotificationChannel{
		ChannelType:              string(models.ChannelTypeMail),
		Id:                       mail.Id,
		ChannelName:              &mail.ChannelName,
		Domain:                   &mail.Domain,
		Port:                     &mail.Port,
		IsAuthenticationRequired: &mail.IsAuthenticationRequired,
		IsTlsEnforced:            &mail.IsTlsEnforced,
		Username:                 mail.Username,
		Password:                 mail.Password,
		MaxEmailAttachmentSizeMb: mail.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    mail.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       &mail.SenderEmailAddress,
	}
}

// MapNotificationChannelsToMail maps a slice of NotificationChannel to MailNotificationChannelRequest.
func MapNotificationChannelsToMail(channels []models.NotificationChannel) []request.MailNotificationChannelRequest {
	mailChannels := make([]request.MailNotificationChannelRequest, 0, len(channels))
	for _, ch := range channels {
		mailChannels = append(mailChannels, MapNotificationChannelToMail(ch))
	}
	return mailChannels
}
