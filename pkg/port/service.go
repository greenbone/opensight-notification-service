// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package port

import (
	"context"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/response"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/dtos"
)

type NotificationService interface {
	ListNotifications(
		ctx context.Context,
		resultSelector query.ResultSelector,
	) (notifications []models.Notification, totalResult uint64, err error)
	CreateNotification(
		ctx context.Context,
		notificationIn models.Notification,
	) (notification models.Notification, err error)
}

type HealthService interface {
	Ready(ctx context.Context) bool
}

type NotificationChannelService interface {
	CreateNotificationChannel(
		ctx context.Context,
		channelIn models.NotificationChannel,
	) (models.NotificationChannel, error)
	GetNotificationChannelByIdAndType(
		ctx context.Context,
		id string,
		channelType models.NotificationChannel,
	) (models.NotificationChannel, error)
	ListNotificationChannelsByType(
		ctx context.Context,
		channelType models.ChannelType,
	) ([]models.NotificationChannel, error)
	UpdateNotificationChannel(
		ctx context.Context,
		id string,
		channelIn models.NotificationChannel,
	) (models.NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id string) error
	CheckNotificationChannelConnectivity(ctx context.Context, channel dtos.CheckMailServerRequest) error
}

type MailChannelService interface {
	CreateMailChannel(
		ctx context.Context,
		channel request.MailNotificationChannelRequest,
	) (request.MailNotificationChannelRequest, error)
}

type MattermostChannelService interface {
	CreateMattermostChannel(
		ctx context.Context,
		channel request.MattermostNotificationChannelRequest,
	) (response.MattermostNotificationChannelResponse, error)
}
