package notificationchannelservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/response"
)

var (
	ErrMattermostChannelLimitReached = errors.New("mattermost channel limit reached")
	ErrListMattermostChannels        = errors.New("failed to list mattermost channels")
	ErrMattermostChannelBadRequest   = errors.New("bad request for mattermost channel")
	ErrMattermostChannelNameExists   = errors.New("mattermost channel name already exists")
)

type MattermostChannelService struct {
	notificationChannelService port.NotificationChannelService
	mattermostChannelLimit     int
}

func NewMattermostChannelService(
	notificationChannelService port.NotificationChannelService,
	mattermostChannelLimit int,
) *MattermostChannelService {
	return &MattermostChannelService{
		notificationChannelService: notificationChannelService,
		mattermostChannelLimit:     mattermostChannelLimit,
	}
}

func (m *MattermostChannelService) mattermostChannelValidations(c context.Context, channelName string) error {
	channels, err := m.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMattermost)
	if err != nil {
		return errors.Join(ErrListMattermostChannels, err)
	}

	if len(channels) >= m.mattermostChannelLimit {
		return ErrMattermostChannelLimitReached
	}

	for _, ch := range channels {
		if ch.ChannelName != nil && *ch.ChannelName == channelName {
			return ErrMattermostChannelNameExists
		}
	}

	return nil
}

func (m *MattermostChannelService) CreateMattermostChannel(
	c context.Context,
	channel request.MattermostNotificationChannelRequest,
) (response.MattermostNotificationChannelResponse, error) {
	if errResp := m.mattermostChannelValidations(c, channel.ChannelName); errResp != nil {
		return response.MattermostNotificationChannelResponse{}, errResp
	}

	notificationChannel := mapper.MapMattermostToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return response.MattermostNotificationChannelResponse{}, err
	}

	return mapper.MapNotificationChannelToMattermost(created), nil
}

func (m *MattermostChannelService) SendMattermostTestMessage(ctx context.Context, id string) error {
	channel, err := m.notificationChannelService.GetNotificationChannelByIdAndType(ctx, id, models.ChannelTypeMattermost)
	if err != nil {
		return err
	}

	if channel.WebhookUrl == nil || *channel.WebhookUrl == "" {
		return ErrMattermostChannelBadRequest
	}

	body, err := json.Marshal(map[string]string{"text": "Hello This is a test message"})
	if err != nil {
		return err
	}

	resp, err := http.Post(*channel.WebhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("failed to send test message to Mattermost webhook: " + resp.Status)
	}

	return nil
}
