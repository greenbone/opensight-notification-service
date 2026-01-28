package notificationchannelservice

import (
	"context"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
)

type NotificationChannelService interface {
	CreateNotificationChannel(
		ctx context.Context,
		channelIn models.NotificationChannel,
	) (models.NotificationChannel, error)
	GetNotificationChannelByIdAndType(
		ctx context.Context,
		id string,
		channelType models.ChannelType,
	) (models.NotificationChannel, error)
	UpdateNotificationChannel(
		ctx context.Context,
		id string,
		channelIn models.NotificationChannel,
	) (models.NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id string) error
	ListNotificationChannelsByType(
		ctx context.Context,
		channelType models.ChannelType,
	) ([]models.NotificationChannel, error)
}
type notificationChannelService struct {
	store notificationrepository.NotificationChannelRepository
}

func NewNotificationChannelService(store notificationrepository.NotificationChannelRepository) NotificationChannelService {
	return &notificationChannelService{
		store: store,
	}
}

func (s *notificationChannelService) CreateNotificationChannel(
	ctx context.Context,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {

	notificationChannel, err := s.store.CreateNotificationChannel(ctx, channelIn)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return notificationChannel, nil
}

func (s *notificationChannelService) GetNotificationChannelByIdAndType(
	ctx context.Context,
	id string,
	channelType models.ChannelType,
) (models.NotificationChannel, error) {
	return s.store.GetNotificationChannelByIdAndType(ctx, id, channelType)
}
func (s *notificationChannelService) ListNotificationChannelsByType(
	ctx context.Context,
	channelType models.ChannelType,
) ([]models.NotificationChannel, error) {
	return s.store.ListNotificationChannelsByType(ctx, channelType)
}

func (s *notificationChannelService) UpdateNotificationChannel(
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

func (s *notificationChannelService) DeleteNotificationChannel(ctx context.Context, id string) error {
	return s.store.DeleteNotificationChannel(ctx, id)
}
