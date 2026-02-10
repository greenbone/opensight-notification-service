package notificationchannelservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/policy"
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
	transport                  http.Client
}

func NewTeamsChannelService(
	notificationChannelService NotificationChannelService,
	teamsChannelLimit int,
	transport http.Client,
) TeamsChannelService {
	return &teamsChannelService{
		notificationChannelService: notificationChannelService,
		teamsChannelLimit:          teamsChannelLimit,
		transport:                  transport,
	}
}

func (t *teamsChannelService) SendTeamsTestMessage(webhookUrl string) error {
	isTeamsOldWebhookUrl, err := policy.IsTeamsOldWebhookUrl(webhookUrl)
	if err != nil {
		// TODO: 10.02.2026 stolksdorf - wrap errors
		return err
	}

	var message map[string]interface{}
	if isTeamsOldWebhookUrl {
		message = map[string]interface{}{
			"text": "Hello, This is a test message",
		}

	} else {
		message = map[string]interface{}{
			"type": "message",
			"attachments": []map[string]interface{}{
				{
					"contentType": "application/vnd.microsoft.card.adaptive",
					"content": map[string]interface{}{
						"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
						"type":    "AdaptiveCard",
						"version": "1.2",
						"body": []map[string]interface{}{
							{
								"type": "TextBlock",
								"text": "Hello, This is a test message",
							},
						},
					},
				},
			},
		}

	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("can not marshal teams message: %w", err)
	}

	resp, err := t.transport.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: http status: %s", ErrTeamsMassageDelivery, resp.Status)
	}

	return nil
}

func (t *teamsChannelService) CreateTeamsChannel(
	ctx context.Context,
	channel teamsdto.TeamsNotificationChannelRequest,
) (teamsdto.TeamsNotificationChannelResponse, error) {
	if err := t.teamsChannelValidations(ctx, channel.ChannelName); err != nil {
		return teamsdto.TeamsNotificationChannelResponse{}, err
	}

	notificationChannel := teamsdto.MapTeamsToNotificationChannel(channel)
	created, err := t.notificationChannelService.CreateNotificationChannel(ctx, notificationChannel)
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
