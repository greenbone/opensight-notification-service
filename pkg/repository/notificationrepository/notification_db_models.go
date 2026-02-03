// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"encoding/json"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
)

const (
	notificationsTable               = "notification_service.notifications"
	createNotificationQuery          = `INSERT INTO ` + notificationsTable + ` (origin, origin_class, origin_resource_id, timestamp, title, detail, level, custom_fields) VALUES (:origin, :origin_class, :origin_resource_id, :timestamp, :title, :detail, :level, :custom_fields) RETURNING *`
	unfilteredListNotificationsQuery = `SELECT * FROM ` + notificationsTable
)

type notificationRow struct {
	Id               string  `db:"id"`
	Origin           string  `db:"origin"`
	OriginClass      string  `db:"origin_class"`
	OriginResourceID *string `db:"origin_resource_id"`
	Timestamp        string  `db:"timestamp"`
	Title            string  `db:"title"`
	Detail           string  `db:"detail"`
	Level            string  `db:"level"`
	CustomFields     []byte  `db:"custom_fields"`
}

func notificationFieldMapping() map[string]string {
	return map[string]string{
		dtos.NameField:            "title",
		dtos.DescriptionFieldName: "detail",
		dtos.LevelFieldName:       "level",
		dtos.OccurrenceFieldName:  "timestamp",
		dtos.OriginFieldName:      "origin",
	}
}

func toNotificationRow(n models.Notification) (notificationRow, error) {
	var empty notificationRow

	customFieldsSerialized, err := json.Marshal(n.CustomFields)
	if err != nil {
		return empty, err
	}

	notificationRow := notificationRow{
		Id:               n.Id,
		Origin:           n.Origin,
		OriginClass:      n.OriginClass,
		OriginResourceID: &n.OriginResourceID,
		Timestamp:        n.Timestamp,
		Title:            n.Title,
		Detail:           n.Detail,
		Level:            n.Level,
		CustomFields:     customFieldsSerialized,
	}

	return notificationRow, nil
}

func (n *notificationRow) ToNotificationModel() (models.Notification, error) {
	var empty models.Notification

	notification := models.Notification{
		Id:               n.Id,
		Origin:           n.Origin,
		OriginClass:      n.OriginClass,
		OriginResourceID: helper.SafeDereference(n.OriginResourceID),
		Timestamp:        n.Timestamp,
		Title:            n.Title,
		Detail:           n.Detail,
		Level:            n.Level,
		// CustomFields is set below
	}

	if len(n.CustomFields) > 0 {
		err := json.Unmarshal(n.CustomFields, &notification.CustomFields)
		if err != nil {
			return empty, err
		}
	}

	return notification, nil
}
