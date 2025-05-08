// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

type Notification struct {
	Id           string         `json:"id" readonly:"true"`
	Origin       string         `json:"origin" binding:"required"`
	OriginUri    string         `json:"originUri,omitempty"` // can be used to provide a link to the origin
	Timestamp    string         `json:"timestamp" binding:"required" format:"date-time"`
	Title        string         `json:"title" binding:"required"`
	Detail       string         `json:"detail" binding:"required"`
	Level        string         `json:"level" binding:"required" enums:"info,warning,error"`
	CustomFields map[string]any `json:"customFields,omitempty"` // can contain arbitrary structured information about the notification
}

// Event is sent by a backend service. It will always result in a notification
// and will possibly also trigger actions like sending mails, depending on the
// cofigured rules.
type Event struct {
	Id               string         `json:"id" readonly:"true"`
	Origin           string         `json:"origin" binding:"required"`      // name of the origin, e.g. `SBOM - React`
	OriginClass      string         `json:"originClass" binding:"required"` // unique identifier for the class of origins, e.g. `/vi/SBOM`
	OriginInstanceID string         `json:"originInstanceID,omitempty"`     // together with class it can be used to provide a link to the origin, e.g. `<id of react sbom object>`
	Timestamp        string         `json:"timestamp" binding:"required" format:"date-time"`
	Title            string         `json:"title" binding:"required"`
	Detail           string         `json:"detail" binding:"required"`
	Level            string         `json:"level" binding:"required" enums:"info,warning,error,urgent"`
	CustomFields     map[string]any `json:"customFields,omitempty"` // can contain arbitrary structured information about the event
}
