// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationservice

import (
	"context"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

type NotificationService struct {
	store port.NotificationRepository
}

func NewNotificationService(store port.NotificationRepository) *NotificationService {
	return &NotificationService{store: store}
}

func (s *NotificationService) ListNotifications(ctx context.Context, resultSelector query.ResultSelector) (notifications []models.Notification, totalResult uint64, err error) {
	return s.store.ListNotifications(ctx, resultSelector)
}

func (s *NotificationService) CreateNotification(ctx context.Context, notificationIn models.Notification) (notification models.Notification, err error) {
	return s.store.CreateNotification(ctx, notificationIn)
}
