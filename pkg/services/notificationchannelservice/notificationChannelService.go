package notificationchannelservice

import (
	"context"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

type NotificationChannelServicer interface {
	CreateNotificationChannel(ctx context.Context, req models.NotificationChannel) (models.NotificationChannel, error)
	GetNotificationChannelByIdAndType(
		ctx context.Context,
		id string,
		channelType models.ChannelType,
	) (models.NotificationChannel, error)
	ListNotificationChannelsByType(
		ctx context.Context,
		channelType models.ChannelType,
	) ([]models.NotificationChannel, error)
	UpdateNotificationChannel(
		ctx context.Context,
		id string,
		req models.NotificationChannel,
	) (models.NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id string) error
	CheckNotificationChannelConnectivity(ctx context.Context, channel models.NotificationChannel) error
	CheckNotificationChannelEntityConnectivity(
		ctx context.Context,
		id string,
		channel models.NotificationChannel,
	) error
}

type NotificationChannelService struct {
	store port.NotificationChannelRepository
}

func NewNotificationChannelService(store port.NotificationChannelRepository) *NotificationChannelService {
	return &NotificationChannelService{store: store}
}

func (s *NotificationChannelService) CreateNotificationChannel(
	ctx context.Context,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {
	notificationChannel, err := s.store.CreateNotificationChannel(ctx, channelIn)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return notificationChannel, nil
}

func (s *NotificationChannelService) GetNotificationChannelByIdAndType(
	ctx context.Context,
	id string,
	channelType models.ChannelType,
) (models.NotificationChannel, error) {
	return s.store.GetNotificationChannelByIdAndType(ctx, id, channelType)
}
func (s *NotificationChannelService) ListNotificationChannelsByType(
	ctx context.Context,
	channelType models.ChannelType,
) ([]models.NotificationChannel, error) {
	return s.store.ListNotificationChannelsByType(ctx, channelType)
}

func (s *NotificationChannelService) UpdateNotificationChannel(
	ctx context.Context,
	id string,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {
	notificationChannel, err := s.store.UpdateNotificationChannel(ctx, id, channelIn)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return notificationChannel, nil
}

func (s *NotificationChannelService) DeleteNotificationChannel(ctx context.Context, id string) error {
	return s.store.DeleteNotificationChannel(ctx, id)
}
