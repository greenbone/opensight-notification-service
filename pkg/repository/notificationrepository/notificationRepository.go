// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"context"
	"errors"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	client *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) (port.NotificationRepository, error) {
	if db == nil {
		return nil, errors.New("nil db reference")
	}
	client := &NotificationRepository{
		client: db,
	}
	return client, nil
}

func (r *NotificationRepository) ListNotifications(ctx context.Context, resultSelector query.ResultSelector) (notifications []models.Notification, totalResult uint64, err error) {
	return
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notificationIn models.Notification) (notification models.Notification, err error) {
	return
}
