package checkmailconnectivity

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
)

const (
	channelListTimeout  = 5 * time.Second
	channelCheckTimeout = 30 * time.Second
)

func NewJob(
	notificationService notificationservice.NotificationService,
	notificationChannelService notificationchannelservice.NotificationChannelService,
	mailChannelService notificationchannelservice.MailChannelService,
) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), channelListTimeout)
		defer cancel()

		mailChannels, err := notificationChannelService.ListNotificationChannelsByType(ctx, models.ChannelTypeMail)
		if err != nil {
			return err
		}

		for _, channel := range mailChannels {
			if err := checkChannelConnectivity(mailChannelService, channel); err != nil {
				_, err := notificationService.CreateNotification(context.Background(), models.Notification{
					Origin:    "Communication service",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Title:     "Mailserver not reachable",
					Detail:    fmt.Sprintf("Mailserver:%s not reachable: %s", *channel.Domain, err),
					Level:     "info",
					CustomFields: map[string]any{
						"Domain":   helper.SafeDereference(channel.Domain),
						"Port":     helper.SafeDereference(channel.Port),
						"Username": helper.SafeDereference(channel.Username),
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

func checkChannelConnectivity(
	mailChannelService notificationchannelservice.MailChannelService,
	channel models.NotificationChannel,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), channelCheckTimeout)
	defer cancel()

	if err := mailChannelService.CheckNotificationChannelConnectivity(ctx, channel); err != nil {
		return err
	}
	return nil
}
