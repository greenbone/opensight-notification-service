package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
)

var (
	ErrMailChannelLimitReached = errors.New("mail channel limit reached")
	ErrListMailChannels        = errors.New("failed to list mail channels")
	ErrGetMailChannel          = errors.New("unable to get notification channel id and type")
)

type MailChannelService struct {
	notificationChannelService port.NotificationChannelService
	emailLimit                 int
}

func NewMailChannelService(
	notificationChannelService port.NotificationChannelService,
	emailLimit int,
) *MailChannelService {
	return &MailChannelService{
		notificationChannelService: notificationChannelService,
		emailLimit:                 emailLimit,
	}
}

func (m *MailChannelService) mailChannelAlreadyExists(c context.Context) error {
	channels, err := m.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		return errors.Join(ErrListMailChannels, err)
	}

	if len(channels) >= m.emailLimit {
		return ErrMailChannelLimitReached
	}
	return nil
}

func (m *MailChannelService) CreateMailChannel(
	c context.Context,
	channel request.MailNotificationChannelRequest,
) (request.MailNotificationChannelRequest, error) {
	if errResp := m.mailChannelAlreadyExists(c); errResp != nil {
		return request.MailNotificationChannelRequest{}, errResp
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return request.MailNotificationChannelRequest{}, err
	}

	return mapper.MapNotificationChannelToMail(created), nil
}

func (s *NotificationChannelService) CheckNotificationChannelConnectivity(
	ctx context.Context,
	mailServer models.NotificationChannel,
) error {
	return ConnectionCheckMail(ctx, mailServer)
}

func (s *NotificationChannelService) CheckNotificationChannelEntityConnectivity(
	ctx context.Context,
	id string,
	mailServer models.NotificationChannel,
) error {
	channel, err := s.GetNotificationChannelByIdAndType(ctx, id, models.ChannelTypeMail)
	if err != nil {
		return errors.Join(ErrGetMailChannel, err)
	}

	if *mailServer.Password == "" && *mailServer.Username != "" {
		if channel.Password != nil {
			mailServer.Password = channel.Password
		}
	}

	return ConnectionCheckMail(ctx, mailServer)
}
