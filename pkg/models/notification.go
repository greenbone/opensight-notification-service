// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/notifications"
	"github.com/greenbone/opensight-notification-service/pkg/validation"
)

// Notification is sent by a backend service. It will always be stored by the
// notification service and it will possibly also trigger actions like sending mails,
// depending on the configured rules.
type Notification struct {
	Id               string              `json:"id" readonly:"true"`
	Origin           string              `json:"origin" validate:"required"` // name of the origin, e.g. `SBOM - React`
	OriginClass      string              `json:"originClass"`                // unique identifier for the class of origins, e.g. `/vi/SBOM`, for now optional for backwards compatibility, will be required in future
	OriginResourceID string              `json:"originResourceID,omitempty"` // together with class it can be used to provide a link to the origin, e.g. `<id of react sbom object>`
	Timestamp        string              `json:"timestamp" validate:"required" format:"date-time"`
	Title            string              `json:"title" validate:"required"`
	Detail           string              `json:"detail" validate:"required"`
	Level            notifications.Level `json:"level" validate:"required" enums:"info,warning,error,urgent"`
	CustomFields     map[string]any      `json:"customFields,omitempty"` // can contain arbitrary structured information about the event
}

func (n *Notification) Validate() ValidationErrors {
	err := validation.Validate.Struct(n)
	if err != nil {
		return ValidationErrors{"$": err.Error()}
	}
	return nil
}
