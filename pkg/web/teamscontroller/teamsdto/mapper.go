package teamsdto

import (
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

// MapNotificationChannelToTeams maps NotificationChannel to TeamsNotificationChannelRequest.
func MapNotificationChannelToTeams(channel models.NotificationChannel) TeamsNotificationChannelResponse {
	return TeamsNotificationChannelResponse{
		Id:          *channel.Id,
		ChannelName: helper.SafeDereference(channel.ChannelName),
		WebhookUrl:  helper.SafeDereference(channel.WebhookUrl),
		Description: helper.SafeDereference(channel.Description),
	}
}

func MapTeamsToNotificationChannel(mail TeamsNotificationChannelRequest) models.NotificationChannel {
	return models.NotificationChannel{
		ChannelType: models.ChannelTypeTeams,
		ChannelName: &mail.ChannelName,
		WebhookUrl:  &mail.WebhookUrl,
		Description: &mail.Description,
	}
}

// MapNotificationChannelsToTeams maps a slice of NotificationChannel to TeamsNotificationChannelRequest.
func MapNotificationChannelsToTeams(channels []models.NotificationChannel) []TeamsNotificationChannelResponse {
	teamsChannels := make([]TeamsNotificationChannelResponse, 0, len(channels))
	for _, ch := range channels {
		teamsChannels = append(teamsChannels, MapNotificationChannelToTeams(ch))
	}
	return teamsChannels
}
