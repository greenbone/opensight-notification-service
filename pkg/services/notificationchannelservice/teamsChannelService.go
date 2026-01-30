package notificationchannelservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamsController/teamsdto"
)

var (
	ErrTeamsChannelLimitReached = errors.New("Teams channel limit reached.")
	ErrListTeamsChannels        = errors.New("failed to list teams channels")
	ErrTeamsChannelNameExists   = errors.New("Teams channel name already exists.")
	ErrTeamsMassageDelivery     = errors.New("teams message could not be send")
)

type TeamsChannelService interface {
	SendTeamsTestMessage(webhookUrl string) error
	CreateTeamsChannel(
		c context.Context,
		channel teamsdto.TeamsNotificationChannelRequest,
	) (teamsdto.TeamsNotificationChannelResponse, error)
}

type teamsChannelService struct {
	notificationChannelService NotificationChannelService
	teamsChannelLimit          int
}

func NewTeamsChannelService(
	notificationChannelService NotificationChannelService,
	teamsChannelLimit int,
) TeamsChannelService {
	return &teamsChannelService{
		notificationChannelService: notificationChannelService,
		teamsChannelLimit:          teamsChannelLimit,
	}
}

func (t *teamsChannelService) SendTeamsTestMessage(webhookUrl string) error {
	body, err := json.Marshal(map[string]string{
		"text": "Hello This is a test message",
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"%w: http status: %s",
			ErrTeamsMassageDelivery,
			resp.Status,
		)
	}

	return nil
}

func (t *teamsChannelService) CreateTeamsChannel(
	c context.Context,
	channel teamsdto.TeamsNotificationChannelRequest,
) (teamsdto.TeamsNotificationChannelResponse, error) {
	if errResp := t.teamsChannelValidations(c, channel.ChannelName); errResp != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, errResp
	}

	notificationChannel := teamsdto.MapTeamsToNotificationChannel(channel)
	created, err := t.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	return teamsdto.MapNotificationChannelToTeams(created), nil
}

func (t *teamsChannelService) teamsChannelValidations(c context.Context, channelName string) error {
	channels, err := t.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeTeams)
	if err != nil {
		return errors.Join(ErrListTeamsChannels, err)
	}

	if len(channels) >= t.teamsChannelLimit {
		return ErrTeamsChannelLimitReached
	}

	for _, ch := range channels {
		if ch.ChannelName != nil && *ch.ChannelName == channelName {
			return ErrTeamsChannelNameExists
		}
	}

	return nil
}
