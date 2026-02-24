package mattermostdto

import (
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

// MapNotificationChannelToMattermost maps NotificationChannel to MattermostNotificationChannelRequest.
func MapNotificationChannelToMattermost(channel models.NotificationChannel) MattermostNotificationChannelResponse {
	return MattermostNotificationChannelResponse{
		Id:          *channel.Id,
		ChannelName: helper.SafeDereference(channel.ChannelName),
		WebhookUrl:  helper.SafeDereference(channel.WebhookUrl),
		Description: helper.SafeDereference(channel.Description),
	}
}

func MapMattermostToNotificationChannel(mail MattermostNotificationChannelRequest) models.NotificationChannel {
	return models.NotificationChannel{
		ChannelType: models.ChannelTypeMattermost,
		ChannelName: &mail.ChannelName,
		WebhookUrl:  &mail.WebhookUrl,
		Description: &mail.Description,
	}
}

// MapNotificationChannelsToMattermost maps a slice of NotificationChannel to MattermostNotificationChannelRequest.
func MapNotificationChannelsToMattermost(channels []models.NotificationChannel) []MattermostNotificationChannelResponse {
	mattermostChannels := make([]MattermostNotificationChannelResponse, 0, len(channels))
	for _, ch := range channels {
		mattermostChannels = append(mattermostChannels, MapNotificationChannelToMattermost(ch))
	}
	return mattermostChannels
}
