package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

var (
	ErrMailChannelAlreadyExists = errors.New("mail channel already exists")
	ErrListMailChannels         = errors.New("failed to list mail channels")
)

type MailChannelService struct {
	notificationChannelService port.NotificationChannelService
}

func NewMailChannelService(notificationChannelService port.NotificationChannelService) *MailChannelService {
	return &MailChannelService{notificationChannelService: notificationChannelService}
}

func (v *MailChannelService) mailChannelAlreadyExists(c context.Context) error {
	channels, err := v.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		return errors.Join(ErrListMailChannels, err)
	}

	if len(channels) > 0 {
		return ErrMailChannelAlreadyExists
	}
	return nil
}

func (v *MailChannelService) CreateMailChannel(c context.Context, channel models.MailNotificationChannel) (models.MailNotificationChannel, error) {
	if errResp := v.mailChannelAlreadyExists(c); errResp != nil {
		return models.MailNotificationChannel{}, errResp
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	created, err := v.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return models.MailNotificationChannel{}, err
	}

	return mapper.MapNotificationChannelToMail(created), nil
}
