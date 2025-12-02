package mapper

import "github.com/greenbone/opensight-notification-service/pkg/models"

// MapNotificationChannelToMail maps NotificationChannel to MailNotificationChannel.
func MapNotificationChannelToMail(channel models.NotificationChannel) models.MailNotificationChannel {
	return models.MailNotificationChannel{
		Id:                       channel.Id,
		ChannelName:              channel.ChannelName,
		Domain:                   channel.Domain,
		Port:                     channel.Port,
		IsAuthenticationRequired: channel.IsAuthenticationRequired,
		IsTlsEnforced:            channel.IsTlsEnforced,
		Username:                 channel.Username,
		Password:                 channel.Password,
		MaxEmailAttachmentSizeMb: channel.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    channel.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       channel.SenderEmailAddress,
	}
}

func MapMailToNotificationChannel(mail models.MailNotificationChannel) models.NotificationChannel {
	return models.NotificationChannel{
		ChannelType:              "mail",
		Id:                       mail.Id,
		ChannelName:              mail.ChannelName,
		Domain:                   mail.Domain,
		Port:                     mail.Port,
		IsAuthenticationRequired: mail.IsAuthenticationRequired,
		IsTlsEnforced:            mail.IsTlsEnforced,
		Username:                 mail.Username,
		MaxEmailAttachmentSizeMb: mail.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    mail.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       mail.SenderEmailAddress,
	}
}

// MapNotificationChannelsToMail maps a slice of NotificationChannel to MailNotificationChannel.
func MapNotificationChannelsToMail(channels []models.NotificationChannel) []models.MailNotificationChannel {
	mailChannels := make([]models.MailNotificationChannel, 0, len(channels))
	for _, ch := range channels {
		mailChannels = append(mailChannels, MapNotificationChannelToMail(ch))
	}
	return mailChannels
}
