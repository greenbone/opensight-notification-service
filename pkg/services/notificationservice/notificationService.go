// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationservice

import (
	"context"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
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

type notificationService struct {
	store notificationrepository.NotificationRepository
}

func NewNotificationService(store notificationrepository.NotificationRepository) NotificationService {
	return &notificationService{store: store}
}

func (s *notificationService) ListNotifications(
	ctx context.Context,
	resultSelector query.ResultSelector,
) (notifications []models.Notification, totalResult uint64, err error) {
	return s.store.ListNotifications(ctx, resultSelector)
}

func (s notificationService) CreateNotification(
	ctx context.Context,
	notificationIn models.Notification,
) (notification models.Notification, err error) {
	return s.store.CreateNotification(ctx, notificationIn)
}
