// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"encoding/json"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

const notificationsTable = "notification_service.notifications"

const createNotificationQuery = `INSERT INTO ` + notificationsTable + ` (origin, origin_uri, timestamp, title, detail, level, custom_fields) VALUES (:origin, :origin_uri, :timestamp, :title, :detail, :level, :custom_fields) RETURNING *`
const unfilteredListNotificationsQuery = `SELECT * FROM ` + notificationsTable

type notificationRow struct {
	Id           string  `db:"id"`
	Origin       string  `db:"origin"`
	OriginUri    *string `db:"origin_uri"`
	Timestamp    string  `db:"timestamp"`
	Title        string  `db:"title"`
	Detail       string  `db:"detail"`
	Level        string  `db:"level"`
	CustomFields []byte  `db:"custom_fields"`
}

func toNotificationRow(n models.Notification) (notificationRow, error) {
	var empty notificationRow

	customFieldsSerialized, err := json.Marshal(n.CustomFields)
	if err != nil {
		return empty, err // TODO: return validation error ?
	}

	notificationRow := notificationRow{
		Id:           n.Id,
		Origin:       n.Origin,
		OriginUri:    &n.OriginUri,
		Timestamp:    n.Timestamp,
		Title:        n.Title,
		Detail:       n.Detail,
		Level:        n.Level,
		CustomFields: customFieldsSerialized,
	}

	return notificationRow, nil
}

func (n *notificationRow) ToNotificationModel() (models.Notification, error) {
	var empty models.Notification

	notification := models.Notification{
		Id:        n.Id,
		Origin:    n.Origin,
		OriginUri: helper.SafeDereference(n.OriginUri),
		Timestamp: n.Timestamp,
		Title:     n.Title,
		Detail:    n.Detail,
		Level:     n.Level,
		// CustomFields is set below
	}

	err := json.Unmarshal(n.CustomFields, &notification.CustomFields)
	if err != nil {
		return empty, err
	}

	return notification, nil
}
