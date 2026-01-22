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
	ErrMattermostChannelLimitReached = errors.New("mattermost channel limit reached")
	ErrListMattermostChannels        = errors.New("failed to list mattermost channels")
	ErrMattermostChannelBadRequest   = errors.New("bad request for mattermost channel")
)

type MattermostChannelService struct {
	notificationChannelService port.NotificationChannelService
	mattermostChannelLimit     int
}

func NewMattermostChannelService(notificationChannelService port.NotificationChannelService, mattermostChannelLimit int) *MattermostChannelService {
	return &MattermostChannelService{notificationChannelService: notificationChannelService, mattermostChannelLimit: mattermostChannelLimit}
}

func (v *MattermostChannelService) mattermostChannelLimitReached(c context.Context) error {
	channels, err := v.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMattermost)
	if err != nil {
		return errors.Join(ErrListMattermostChannels, err)
	}

	if len(channels) >= v.mattermostChannelLimit {
		return ErrMattermostChannelLimitReached
	}
	return nil
}

func (v *MattermostChannelService) CreateMattermostChannel(c context.Context, channel request.MattermostNotificationChannelRequest) (response.MattermostNotificationChannelResponse, error) {
	if errResp := v.mattermostChannelLimitReached(c); errResp != nil {
		return response.MattermostNotificationChannelResponse{}, errResp
	}

	notificationChannel := mapper.MapMattermostToNotificationChannel(channel)
	created, err := v.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return response.MattermostNotificationChannelResponse{}, err
	}

	return mapper.MapNotificationChannelToMattermost(created), nil
}
