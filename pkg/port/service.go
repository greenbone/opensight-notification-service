// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package port

import (
	"context"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

type NotificationService interface {
	ListNotifications(ctx context.Context, resultSelector query.ResultSelector) (notifications []models.Notification, totalResult uint64, err error)
	CreateNotification(ctx context.Context, notificationIn models.Notification) (notification models.Notification, err error)
}

type HealthService interface {
	Ready(ctx context.Context) bool
}

type NotificationChannelService interface {
	CreateNotificationChannel(ctx context.Context, channelIn models.NotificationChannel) (models.NotificationChannel, error)
	ListNotificationChannelsByType(ctx context.Context, channelType string) ([]models.NotificationChannel, error)
	UpdateNotificationChannel(ctx context.Context, id string, channelIn models.NotificationChannel) (models.NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id string) error
}
