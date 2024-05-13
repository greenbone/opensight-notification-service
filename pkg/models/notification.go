// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

type Notification struct {
	Id           string         `json:"id" readonly:"true"`
	Origin       string         `json:"origin" binding:"required"`
	OriginUri    string         `json:"originUri,omitempty"` // can be used to provide a link to the origin
	Timestamp    string         `json:"timestamp" binding:"required" format:"date-time"`
	Title        string         `json:"title" binding:"required"` // can also be seen as the 'type'
	Detail       string         `json:"detail" binding:"required"`
	Level        string         `json:"level" binding:"required" enums:"info,warning,error,critical"`
	CustomFields map[string]any `json:"customFields,omitempty"` // can contain arbitrary structured information about the notification
}
