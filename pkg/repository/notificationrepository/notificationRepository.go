// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"context"
	"errors"
	"fmt"

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

func (r *NotificationRepository) ListNotifications(ctx context.Context, resultSelector query.ResultSelector) (notifications []models.Notification, totalResults uint64, err error) {
	var rows []notificationRow
	query := unfilteredListNotificationsQuery

	if resultSelector.Paging != nil { // TODO: add support for filtering and sorting
		limit := resultSelector.Paging.PageSize
		offset := resultSelector.Paging.PageIndex * resultSelector.Paging.PageSize
		query += fmt.Sprint(` LIMIT `, limit, ` OFFSET `, offset)
	}

	err = r.client.SelectContext(ctx, &rows, query)
	if err != nil {
		err = fmt.Errorf("error getting notifications from database: %w", err)
		return
	}

	err = r.client.QueryRowxContext(ctx, `SELECT count(*) FROM `+notificationsTable).Scan(&totalResults)
	if err != nil {
		err = fmt.Errorf("error getting total results: %w", err)
		return
	}

	notifications = make([]models.Notification, 0, len(rows))
	for _, row := range rows {
		notification, err := row.ToNotificationModel()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to transform notification db entry: %w", err)
		}
		notifications = append(notifications, notification)
	}
	return
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notificationIn models.Notification) (notification models.Notification, err error) {
	insertRow, err := toNotificationRow(notificationIn)
	if err != nil {
		return notification, fmt.Errorf("invalid argument for inserting notification into database: %w", err)
	}

	createNotificationStatement, err := r.client.PrepareNamedContext(ctx, createNotificationQuery)
	if err != nil {
		return notification, fmt.Errorf("could not prepare sql statement: %w", err)
	}

	var row notificationRow
	err = createNotificationStatement.QueryRowxContext(ctx, insertRow).StructScan(&row)
	if err != nil {
		return notification, fmt.Errorf("could not insert into database: %w", err)
	}

	notification, err = row.ToNotificationModel()
	if err != nil {
		return notification, fmt.Errorf("failed to transform notification db entry to model: %w", err)
	}

	return notification, nil
}
