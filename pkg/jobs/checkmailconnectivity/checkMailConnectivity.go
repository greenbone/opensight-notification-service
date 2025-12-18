package checkmailconnectivity

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

const (
	channelListTimeout  = 5 * time.Second
	channelCheckTimeout = 30 * time.Second
)

func NewJob(service port.NotificationChannelService) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), channelListTimeout)
		defer cancel()
		li, err := service.ListNotificationChannelsByType(ctx, "mail")
		if err != nil {
			return err
		}
		for _, channel := range li {
			if err := checkChannelConnectivity(service, channel); err != nil {
				return fmt.Errorf("channel %q: %w", *channel.Id, err)
			}
		}
		return nil
	}
}

func checkChannelConnectivity(service port.NotificationChannelService, channel models.NotificationChannel) error {
	ctx, cancel := context.WithTimeout(context.Background(), channelCheckTimeout)
	defer cancel()
	if err := service.CheckNotificationChannelConnectivity(ctx, channel); err != nil {
		return err // TODO: notify about failed check instead of returning error
	}
	return nil
}
