// SPDX-FileCopyrightText: 2023 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

type Notification struct {
	Id           string         `json:"id" db:"id" readonly:"true"`
	Origin       string         `json:"origin" db:"origin" binding:"required"`
	OriginUri    string         `json:"originUri,omitempty" db:"origin"` // can be used to provide a link to the origin
	Timestamp    string         `json:"timestamp" db:"timestamp"  binding:"required" format:"date-time"`
	Title        string         `json:"title" db:"title" binding:"required"` // can also be seen as the 'type'
	Detail       string         `json:"detail" db:"detail" binding:"required"`
	Level        string         `json:"level" db:"level" binding:"required" enums:"info,warning,error,critical"`
	CustomFields map[string]any `json:"customFields,omitempty" db:"custom_fields"` // can contain arbitrary structured information about the notification
}
