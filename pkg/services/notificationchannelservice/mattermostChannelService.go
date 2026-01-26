package notificationchannelservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/dto"
)

var (
	ErrMattermostChannelLimitReached = errors.New("mattermost channel limit reached")
	ErrListMattermostChannels        = errors.New("failed to list mattermost channels")
	ErrMattermostChannelNameExists   = errors.New("mattermost channel name already exists")
)

type MattermostChannelService interface {
	SendMattermostTestMessage(webhookUrl string) error
	CreateMattermostChannel(
		c context.Context,
		channel dto.MattermostNotificationChannelRequest,
	) (dto.MattermostNotificationChannelResponse, error)
}

type mattermostChannelService struct {
	notificationChannelService NotificationChannelService
	mattermostChannelLimit     int
}

func (m *mattermostChannelService) SendMattermostTestMessage(webhookUrl string) error {
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
		return fmt.Errorf("failed to send test message to Mattermost webhook: %s", resp.Status)
	}

	return nil
}

func (m *mattermostChannelService) CreateMattermostChannel(
	c context.Context,
	channel dto.MattermostNotificationChannelRequest,
) (dto.MattermostNotificationChannelResponse, error) {
	if errResp := m.mattermostChannelValidations(c, channel.ChannelName); errResp != nil {
		return dto.MattermostNotificationChannelResponse{}, errResp
	}

	notificationChannel := dto.MapMattermostToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return dto.MattermostNotificationChannelResponse{}, err
	}

	return dto.MapNotificationChannelToMattermost(created), nil
}

func NewMattermostChannelService(
	notificationChannelService NotificationChannelService,
	mattermostChannelLimit int,
) MattermostChannelService {
	return &mattermostChannelService{
		notificationChannelService: notificationChannelService,
		mattermostChannelLimit:     mattermostChannelLimit,
	}
}

func (m *mattermostChannelService) mattermostChannelValidations(c context.Context, channelName string) error {
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
