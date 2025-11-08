package notificationchannelservice

import (
	"context"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

type NotificationChannelService struct {
	store port.NotificationChannelRepository
}

func NewNotificationChannelService(store port.NotificationChannelRepository) *NotificationChannelService {
	return &NotificationChannelService{store: store}
}

func (s *NotificationChannelService) CreateNotificationChannel(ctx context.Context, channelIn models.NotificationChannel) (models.NotificationChannel, error) {
	return s.store.CreateNotificationChannel(ctx, channelIn)
}

func (s *NotificationChannelService) ListNotificationChannelsByType(ctx context.Context, channelType string) ([]models.NotificationChannel, error) {
	return s.store.ListNotificationChannelsByType(ctx, channelType)
}

func (s *NotificationChannelService) UpdateNotificationChannel(ctx context.Context, id string, channelIn models.NotificationChannel) (models.NotificationChannel, error) {
	return s.store.UpdateNotificationChannel(ctx, id, channelIn)
}

func (s *NotificationChannelService) DeleteNotificationChannel(ctx context.Context, id string) error {
	return s.store.DeleteNotificationChannel(ctx, id)
}
