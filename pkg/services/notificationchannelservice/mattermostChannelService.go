package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/mattermostdto"
)

var (
	ErrMattermostChannelLimitReached = errors.New("Mattermost channel limit reached.")
	ErrListMattermostChannels        = errors.New("failed to list mattermost channels")
	ErrMattermostChannelNameExists   = errors.New("Mattermost channel name already exists.")
	ErrMattermostMassageDelivery     = errors.New("mattermost message could not be send")
)

type MattermostChannelService interface {
	SendMattermostTestMessage(webhookUrl string) error
	CreateMattermostChannel(
		ctx context.Context,
		channel mattermostdto.MattermostNotificationChannelRequest,
	) (mattermostdto.MattermostNotificationChannelResponse, error)
	UpdateMattermostChannel(
		ctx context.Context,
		id string,
		channel mattermostdto.MattermostNotificationChannelRequest,
	) (mattermostdto.MattermostNotificationChannelResponse, error)
}

type mattermostChannelService struct {
	notificationChannelService NotificationChannelService
	mattermostChannelLimit     int
	mattermostService          *MattermostService
}

func NewMattermostChannelService(
	notificationChannelService NotificationChannelService,
	mattermostChannelLimit int,
	mattermostService *MattermostService,
) MattermostChannelService {
	return &mattermostChannelService{
		notificationChannelService: notificationChannelService,
		mattermostChannelLimit:     mattermostChannelLimit,
		mattermostService:          mattermostService,
	}
}

func (m *mattermostChannelService) SendMattermostTestMessage(webhookUrl string) error {
	return m.mattermostService.SendMessage(webhookUrl, "Hello, This is a test message")
}

func (m *mattermostChannelService) CreateMattermostChannel(
	ctx context.Context,
	channel mattermostdto.MattermostNotificationChannelRequest,
) (mattermostdto.MattermostNotificationChannelResponse, error) {
	if err := m.mattermostChannelValidations(ctx, channel.ChannelName, ""); err != nil {
		return mattermostdto.MattermostNotificationChannelResponse{}, err
	}

	notificationChannel := mattermostdto.MapMattermostToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return mattermostdto.MattermostNotificationChannelResponse{}, err
	}

	return mattermostdto.MapNotificationChannelToMattermost(created), nil
}

func (m *mattermostChannelService) UpdateMattermostChannel(
	ctx context.Context,
	id string,
	channel mattermostdto.MattermostNotificationChannelRequest,
) (mattermostdto.MattermostNotificationChannelResponse, error) {

	if err := m.mattermostChannelValidations(ctx, channel.ChannelName, id); err != nil {
		return mattermostdto.MattermostNotificationChannelResponse{}, err
	}

	notificationChannel := mattermostdto.MapMattermostToNotificationChannel(channel)
	updated, err := m.notificationChannelService.UpdateNotificationChannel(ctx, id, notificationChannel)
	if err != nil {
		return mattermostdto.MattermostNotificationChannelResponse{}, err
	}

	return mattermostdto.MapNotificationChannelToMattermost(updated), nil
}

func (m *mattermostChannelService) mattermostChannelValidations(
	ctx context.Context,
	channelName string,
	excludeId string,
) error {
	channels, err := m.notificationChannelService.ListNotificationChannelsByType(ctx, models.ChannelTypeMattermost)
	if err != nil {
		return errors.Join(ErrListMattermostChannels, err)
	}

	if len(channels) >= m.mattermostChannelLimit {
		return ErrMattermostChannelLimitReached
	}

	for _, ch := range channels {
		if ch.Id == excludeId {
			continue
		}

		if ch.ChannelName == channelName {
			return ErrMattermostChannelNameExists
		}
	}

	return nil
}
