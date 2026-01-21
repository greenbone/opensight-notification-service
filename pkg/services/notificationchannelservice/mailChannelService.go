package notificationchannelservice

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/dtos"
)

var (
	ErrMailChannelAlreadyExists = errors.New("mail channel already exists")
	ErrListMailChannels         = errors.New("failed to list mail channels")
	ErrMailChannelBadRequest    = errors.New("bad request for mail channel")
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

func (v *MailChannelService) CreateMailChannel(c context.Context, channel request.MailNotificationChannelRequest) (request.MailNotificationChannelRequest, error) {
	if errResp := v.mailChannelAlreadyExists(c); errResp != nil {
		return request.MailNotificationChannelRequest{}, errResp
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)
	created, err := v.notificationChannelService.CreateNotificationChannel(c, notificationChannel)
	if err != nil {
		return request.MailNotificationChannelRequest{}, err
	}

	return mapper.MapNotificationChannelToMail(created), nil
}

func (s *NotificationChannelService) CheckNotificationChannelConnectivity(
	ctx context.Context,
	mailServer dtos.CheckMailServerRequest,
) error {
	return ConnectionCheckMail(ctx, mailServer)
}
