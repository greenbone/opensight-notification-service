package notificationchannelservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/rs/zerolog/log"
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
	store          port.NotificationChannelRepository
	encryptManager port.EncryptManager
}

func NewNotificationChannelService(store port.NotificationChannelRepository, encryptManager port.EncryptManager) *NotificationChannelService {
	return &NotificationChannelService{
		store:          store,
		encryptManager: encryptManager,
	}
}

func (s *NotificationChannelService) CreateNotificationChannel(
	ctx context.Context,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {

	channelInWithEncryptedValues, err := s.encryptValues(channelIn)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	notificationChannel, err := s.store.CreateNotificationChannel(ctx, channelInWithEncryptedValues)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	channelInWithDecryptedValues, err := s.decryptValues(notificationChannel)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return channelInWithDecryptedValues, nil
}

func (s *NotificationChannelService) GetNotificationChannelByIdAndType(
	ctx context.Context,
	id string,
	channelType models.ChannelType,
) (models.NotificationChannel, error) {

	notificationChannel, err := s.store.GetNotificationChannelByIdAndType(ctx, id, channelType)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	channelInWithDecryptedValues, err := s.decryptValues(notificationChannel)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return channelInWithDecryptedValues, nil
}
func (s *NotificationChannelService) ListNotificationChannelsByType(
	ctx context.Context,
	channelType models.ChannelType,
) ([]models.NotificationChannel, error) {

	var decryptedChannels []models.NotificationChannel

	channels, err := s.store.ListNotificationChannelsByType(ctx, channelType)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		decryptedChannel, err := s.decryptValues(channel)
		if err != nil {
			return nil, err
		}

		decryptedChannels = append(decryptedChannels, decryptedChannel)
	}

	return decryptedChannels, nil
}

func (s *NotificationChannelService) UpdateNotificationChannel(
	ctx context.Context,
	id string,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {

	channelInWithEncryptedValues, err := s.encryptValues(channelIn)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	notificationChannel, err := s.store.UpdateNotificationChannel(ctx, id, channelInWithEncryptedValues)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	channelInWithDecryptedValues, err := s.decryptValues(notificationChannel)
	if err != nil {
		return models.NotificationChannel{}, err
	}

	return channelInWithDecryptedValues, nil
}

func (s *NotificationChannelService) DeleteNotificationChannel(ctx context.Context, id string) error {
	return s.store.DeleteNotificationChannel(ctx, id)
}

func (s *NotificationChannelService) encryptValues(row models.NotificationChannel) (models.NotificationChannel, error) {
	if row.Password != nil && strings.TrimSpace(*row.Password) != "" {
		encryptedPasswd, err := s.encryptManager.Encrypt(*row.Password)
		if err != nil {
			return models.NotificationChannel{}, fmt.Errorf("could not encrypt password: %w", err)
		}

		passwd := string(encryptedPasswd)
		row.Password = &passwd
	}

	if row.Username != nil && strings.TrimSpace(*row.Username) != "" {
		encryptedUsername, err := s.encryptManager.Encrypt(*row.Username)
		if err != nil {
			return models.NotificationChannel{}, fmt.Errorf("could not encrypt password: %w", err)
		}

		username := string(encryptedUsername)
		row.Username = &username
	}

	return row, nil
}

func (s *NotificationChannelService) decryptValues(row models.NotificationChannel) (models.NotificationChannel, error) {
	if row.Password != nil && strings.TrimSpace(*row.Password) != "" {
		dPasswd := *row.Password
		dcPassword, err := s.encryptManager.Decrypt([]byte(dPasswd))
		if err != nil {
			log.Err(err).Msg("could not decrypt password")
		}

		row.Password = &dcPassword
	}

	if row.Username != nil && strings.TrimSpace(*row.Username) != "" {
		username := *row.Username
		dcUsername, err := s.encryptManager.Decrypt([]byte(username))
		if err != nil {
			log.Err(err).Msg("could not decrypt password")
		}

		row.Username = &dcUsername
	}

	return row, nil
}
