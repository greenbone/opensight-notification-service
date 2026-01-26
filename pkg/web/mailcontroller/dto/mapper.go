package dto

import "github.com/greenbone/opensight-notification-service/pkg/models"

// MapNotificationChannelToMail maps NotificationChannel to MailNotificationChannelRequest.
func MapNotificationChannelToMail(channel models.NotificationChannel) MailNotificationChannelRequest {
	return MailNotificationChannelRequest{
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

func MapMailToNotificationChannel(mail MailNotificationChannelRequest) models.NotificationChannel {
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

// MapNotificationChannelsToMailWithEmptyPassword maps a slice of NotificationChannel to MailNotificationChannelRequest.
func MapNotificationChannelsToMailWithEmptyPassword(channels []models.NotificationChannel) []MailNotificationChannelRequest {
	mailChannels := make([]MailNotificationChannelRequest, 0, len(channels))
	for _, ch := range channels {
		mailChannels = append(mailChannels, MapNotificationChannelToMail(ch).WithEmptyPassword())
	}
	return mailChannels
}
