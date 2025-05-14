// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

// Notification is sent by a backend service. It will always be stored by the
// notification service and it will possibly also trigger actions like sending mails,
// depending on the cofigured rules.
type Notification struct {
	Id               string         `json:"id" readonly:"true"`
	Origin           string         `json:"origin" binding:"required"`      // name of the origin, e.g. `SBOM - React`
	OriginClass      string         `json:"originClass" binding:"required"` // unique identifier for the class of origins, e.g. `/vi/SBOM`
	OriginResourceID string         `json:"originResourceID,omitempty"`     // together with class it can be used to provide a link to the origin, e.g. `<id of react sbom object>`
	Timestamp        string         `json:"timestamp" binding:"required" format:"date-time"`
	Title            string         `json:"title" binding:"required"`
	Detail           string         `json:"detail" binding:"required"`
	Level            string         `json:"level" binding:"required" enums:"info,warning,error,urgent"`
	CustomFields     map[string]any `json:"customFields,omitempty"` // can contain arbitrary structured information about the event
}
