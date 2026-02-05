// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import "github.com/greenbone/opensight-notification-service/pkg/entities"

// Origin of an event/notification, uniquely identified by the combination of namespace and class.
type Origin struct {
	Name      string `json:"name" binding:"required"`
	Class     string `json:"class" binding:"required"`
	Namespace string `json:"namespace" readonly:"true"`
}

// ToEntity transforms the rest model to the entity for use in the service
func (o Origin) ToEntity() entities.Origin {
	return entities.Origin(o)
}
