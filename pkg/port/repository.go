// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package port

import (
	"context"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

type NotificationRepository interface {
	ListNotifications(
		ctx context.Context,
		resultSelector query.ResultSelector,
	) (notifications []models.Notification, totalResult uint64, err error)
	CreateNotification(
		ctx context.Context,
		notificationIn models.Notification,
	) (notification models.Notification, err error)
}

type NotificationChannelRepository interface {
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
	DeleteNotificationChannel(ctx context.Context, id string) error
	UpdateNotificationChannel(
		ctx context.Context,
		id string,
		in models.NotificationChannel,
	) (models.NotificationChannel, error)
}
