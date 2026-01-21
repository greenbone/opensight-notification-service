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

// MapNotificationChannelToMattermost maps NotificationChannel to MattermostNotificationChannelRequest.
func MapNotificationChannelToMattermost(channel models.NotificationChannel) request.MattermostNotificationChannelRequest {
	return request.MattermostNotificationChannelRequest{
		Id:          channel.Id,
		ChannelName: *channel.ChannelName,
		WebhookUrl:  *channel.WebhookUrl,
		Description: *channel.Description,
	}
}

func MapMattermostToNotificationChannel(mail request.MattermostNotificationChannelRequest) models.NotificationChannel {
	return models.NotificationChannel{
		ChannelType: string(models.ChannelTypeMattermost),
		Id:          mail.Id,
		ChannelName: &mail.ChannelName,
		WebhookUrl:  &mail.WebhookUrl,
		Description: &mail.Description,
	}
}

// MapNotificationChannelsToMattermost maps a slice of NotificationChannel to MattermostNotificationChannelRequest.
func MapNotificationChannelsToMattermost(channels []models.NotificationChannel) []request.MattermostNotificationChannelRequest {
	mattermostChannels := make([]request.MattermostNotificationChannelRequest, 0, len(channels))
	for _, ch := range channels {
		mattermostChannels = append(mattermostChannels, MapNotificationChannelToMattermost(ch))
	}
	return mattermostChannels
}
