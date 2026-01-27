package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/maildto"
)

var (
	ErrMailChannelLimitReached = errors.New("mail channel limit reached")
	ErrListMailChannels        = errors.New("failed to list mail channels")
)

type MailChannelService interface {
	CreateMailChannel(
		c context.Context,
		channel maildto.MailNotificationChannelRequest,
	) (maildto.MailNotificationChannelResponse, error)
	CheckNotificationChannelConnectivity(
		ctx context.Context,
		mailServer models.NotificationChannel,
	) error
	CheckNotificationChannelEntityConnectivity(
		ctx context.Context,
		id string,
		mailServer models.NotificationChannel,
	) error
}

type mailChannelService struct {
	notificationChannelService NotificationChannelService
	store                      notificationrepository.NotificationChannelRepository
	emailLimit                 int
}

func NewMailChannelService(
	notificationChannelService NotificationChannelService,
	store notificationrepository.NotificationChannelRepository,
	emailLimit int,
) MailChannelService {
	return &mailChannelService{
		notificationChannelService: notificationChannelService,
		store:                      store,
		emailLimit:                 emailLimit,
	}
}

func (m *mailChannelService) CreateMailChannel(
	c context.Context,
	channel maildto.MailNotificationChannelRequest,
) (maildto.MailNotificationChannelResponse, error) {
	if errResp := m.mailChannelAlreadyExists(c); errResp != nil {
		return maildto.MailNotificationChannelResponse{}, errResp
	}

	notificationChannel := maildto.MapMailToNotificationChannel(channel)
	created, err := m.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return maildto.MailNotificationChannelResponse{}, err
	}

	return maildto.MapNotificationChannelToMail(created), nil
}

func (m *mailChannelService) CheckNotificationChannelConnectivity(
	ctx context.Context,
	mailServer models.NotificationChannel,
) error {
	return ConnectionCheckMail(ctx, mailServer)
}

func (m *mailChannelService) CheckNotificationChannelEntityConnectivity(
	ctx context.Context,
	id string,
	mailServer models.NotificationChannel,
) error {
	channel, err := m.store.GetNotificationChannelByIdAndType(ctx, id, models.ChannelTypeMail)
	if err != nil {
		return err
	}

	if *mailServer.Password == "" && *mailServer.Username != "" {
		if channel.Password != nil {
			mailServer.Password = channel.Password
		}
	}

	return ConnectionCheckMail(ctx, mailServer)
}

func (m *mailChannelService) mailChannelAlreadyExists(c context.Context) error {
	channels, err := m.notificationChannelService.ListNotificationChannelsByType(c, models.ChannelTypeMail)
	if err != nil {
		return errors.Join(ErrListMailChannels, err)
	}

	if len(channels) >= m.emailLimit {
		return ErrMailChannelLimitReached
	}
	return nil
}
