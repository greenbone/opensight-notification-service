package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/response"
)

var (
	ErrTeamsChannelLimitReached = errors.New("teams channel limit reached")
	ErrListTeamsChannels        = errors.New("failed to list teams channels")
	ErrTeamsChannelBadRequest   = errors.New("bad request for teams channel")
)

type TeamsChannelService struct {
	notificationChannelService port.NotificationChannelService
	teamsChannelLimit          int
}

func NewTeamsChannelService(
	notificationChannelService port.NotificationChannelService,
	teamsChannelLimit int,
) *TeamsChannelService {
	return &TeamsChannelService{
		notificationChannelService: notificationChannelService,
		teamsChannelLimit:          teamsChannelLimit,
	}
}

func (m *TeamsChannelService) teamsChannelLimitReached(c context.Context) error {
	channels, err := m.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeTeams)
	if err != nil {
		return errors.Join(ErrListTeamsChannels, err)
	}

	if len(channels) >= m.teamsChannelLimit {
		return ErrTeamsChannelLimitReached
	}
	return nil
}

func (m *TeamsChannelService) CreateTeamsChannel(
	c context.Context,
	channel request.TeamsNotificationChannelRequest,
) (response.TeamsNotificationChannelResponse, error) {
	if errResp := m.teamsChannelLimitReached(c); errResp != nil {
		return response.TeamsNotificationChannelResponse{}, errResp
	}

	notificationChannel := mapper.MapTeamsToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return response.TeamsNotificationChannelResponse{}, err
	}

	return mapper.MapNotificationChannelToTeams(created), nil
}
