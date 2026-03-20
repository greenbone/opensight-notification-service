package notificationchannelservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamscontroller/teamsdto"
)

var (
	ErrTeamsChannelLimitReached = errors.New("Teams channel limit reached.")
	ErrListTeamsChannels        = errors.New("failed to list teams channels")
	ErrTeamsChannelNameExists   = errors.New("Teams channel name already exists.")
	ErrTeamsMessageDelivery     = errors.New("teams message could not be send")
)

type TeamsChannelService interface {
	SendTeamsTestMessage(webhookUrl string) error
	CreateTeamsChannel(
		ctx context.Context,
		channel teamsdto.TeamsNotificationChannelRequest,
	) (teamsdto.TeamsNotificationChannelResponse, error)
	UpdateTeamsChannel(
		ctx context.Context,
		id string,
		channel teamsdto.TeamsNotificationChannelRequest,
	) (teamsdto.TeamsNotificationChannelResponse, error)
}

type teamsChannelService struct {
	notificationChannelService NotificationChannelService
	teamsChannelLimit          int
	transport                  *http.Client
}

func NewTeamsChannelService(
	notificationChannelService NotificationChannelService,
	teamsChannelLimit int,
	transport *http.Client,
) TeamsChannelService {
	return &teamsChannelService{
		notificationChannelService: notificationChannelService,
		teamsChannelLimit:          teamsChannelLimit,
		transport:                  transport,
	}
}

func (t *teamsChannelService) SendTeamsTestMessage(webhookUrl string) error {
	return sendTeamsMessage(t.transport, webhookUrl, "Hello, This is a test message")
}

func (t *teamsChannelService) CreateTeamsChannel(
	ctx context.Context,
	channel teamsdto.TeamsNotificationChannelRequest,
) (teamsdto.TeamsNotificationChannelResponse, error) {
	if err := t.teamsChannelValidations(ctx, channel.ChannelName, ""); err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	notificationChannel := teamsdto.MapTeamsToNotificationChannel(channel)
	created, err := t.notificationChannelService.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	return teamsdto.MapNotificationChannelToTeams(created), nil
}

func (t *teamsChannelService) UpdateTeamsChannel(
	ctx context.Context,
	id string,
	channel teamsdto.TeamsNotificationChannelRequest,
) (teamsdto.TeamsNotificationChannelResponse, error) {
	if err := t.teamsChannelValidations(ctx, channel.ChannelName, id); err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	notificationChannel := teamsdto.MapTeamsToNotificationChannel(channel)
	created, err := t.notificationChannelService.UpdateNotificationChannel(ctx, id, notificationChannel)
	if err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	return teamsdto.MapNotificationChannelToTeams(created), nil
}

func (t *teamsChannelService) teamsChannelValidations(
	ctx context.Context,
	channelName string,
	excludeId string,
) error {
	channels, err := t.notificationChannelService.ListNotificationChannelsByType(ctx, models.ChannelTypeTeams)
	if err != nil {
		return errors.Join(ErrListTeamsChannels, err)
	}

	if len(channels) >= t.teamsChannelLimit {
		return ErrTeamsChannelLimitReached
	}

	for _, ch := range channels {
		if ch.Id == excludeId {
			continue
		}

		if ch.ChannelName == channelName {
			return ErrTeamsChannelNameExists
		}
	}

	return nil
}
