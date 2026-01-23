package checkmailconnectivity

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
)

const (
	channelListTimeout  = 5 * time.Second
	channelCheckTimeout = 30 * time.Second
)

func NewJob(
	notificationService *notificationservice.NotificationService,
	service port.NotificationChannelService,
) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), channelListTimeout)
		defer cancel()

		mailChannels, err := service.ListNotificationChannelsByType(ctx, models.ChannelTypeMail)
		if err != nil {
			return err
		}

		for _, channel := range mailChannels {
			if err := checkChannelConnectivity(service, channel); err != nil {
				_, err := notificationService.CreateNotification(context.Background(), models.Notification{
					Origin:    "Communication service",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Title:     "Mailserver not reachable",
					Detail:    fmt.Sprintf("Mailserver:%s not reachable: %s", *channel.Domain, err),
					Level:     "info",
					CustomFields: map[string]any{
						"Domain":   Value(channel.Domain),
						"Port":     Value(channel.Port),
						"Username": Value(channel.Username),
					},
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func Value[T any](value *T) any {
	if value == nil {
		return nil
	}
	return *value
}

func checkChannelConnectivity(service port.NotificationChannelService, channel models.NotificationChannel) error {
	ctx, cancel := context.WithTimeout(context.Background(), channelCheckTimeout)
	defer cancel()

	if err := service.CheckNotificationChannelConnectivity(ctx, channel); err != nil {
		return err
	}
	return nil
}
