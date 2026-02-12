// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import (
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/validation"
)

// Origin of an event/notification.
type Origin struct {
	Name      string `json:"name" validate:"required"`  // human readable name representation
	Class     string `json:"class" validate:"required"` // unique identifier
	ServiceID string `json:"serviceID" readonly:"true"` // service in which this origin is defined
}

// ToEntity transforms the rest model to the entity for use in the service
func (o Origin) ToEntity() entities.Origin {
	return entities.Origin(o)
}

type OriginList []Origin

func (o OriginList) Validate() ValidationErrors {
	err := validation.Validate.Var(o, "omitempty,dive")
	if err != nil {
		return ValidationErrors{"$": err.Error()}
	}
	return nil
}
