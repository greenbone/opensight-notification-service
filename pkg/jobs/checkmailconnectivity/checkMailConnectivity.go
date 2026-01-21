package checkmailconnectivity

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/dtos"
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

		mailChannels, err := service.ListNotificationChannelsByType(ctx, "mail")
		if err != nil {
			return err
		}

		for _, channel := range mailChannels {
			// TODO: 21.01.2026 stolksdorf - fix asd
			asd := dtos.NewCheckMailServerRequest(channel)
			if err := checkChannelConnectivity(service, asd); err != nil {
				_, err := notificationService.CreateNotification(ctx, models.Notification{
					Origin:    "notification-service",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Title:     "Mail server not reachable",
					Detail:    fmt.Sprintf("Mail server:%s not reachable: %s", *channel.Domain, err),
					Level:     "info",
					CustomFields: map[string]any{
						"Domain":   *channel.Domain,
						"Port":     *channel.Port,
						"Username": *channel.Username,
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

func checkChannelConnectivity(service port.NotificationChannelService, channel dtos.CheckMailServerRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), channelCheckTimeout)
	defer cancel()

	if err := service.CheckNotificationChannelConnectivity(ctx, channel); err != nil {
		return err
	}
	return nil
}
