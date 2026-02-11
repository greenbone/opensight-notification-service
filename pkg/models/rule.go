// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import (
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/validation"
)

// A rule determines which events cause which action.
// Each incoming event is matched with the trigger conditions.
// If the condition is fulfilled, the provided action is triggered.
type Rule struct {
	ID      string  `json:"id" readonly:"true"`
	Name    string  `json:"name" validate:"required"`
	Trigger Trigger `json:"trigger" validate:"required"`
	Action  Action  `json:"action" validate:"required"`
	Active  bool    `json:"active"`
}

// Trigger condition, fulfilled if both one of `origins` and `levels` match the ones from the incomming event.
type Trigger struct {
	Origins []OriginReference `json:"origins" validate:"required"`
	Levels  []string          `json:"level" validate:"required"`
}

// Action determines to which channel the event is forwarded.
// Some channels (e.g. mail) require the explicit recipient(s).
type Action struct {
	Channel   ChannelReference `json:"channel" validate:"required"`
	Recipient string           `json:"recipient,omitempty"` // specific recipient if supported/required by the channel, e.g. for mail a comma separated list of mail adresses
}

type RuleOptions struct {
	Channels     []ChannelReference `json:"channels"`
	EventOrigins []OriginReference  `json:"eventOrigins"`
	EventLevels  []string           `json:"eventLevels"`
}

type OriginReference struct {
	Name      string `json:"name" readonly:"true"`
	Class     string `json:"class" validate:"required"`
	ServiceID string `json:"serviceID" validate:"required"`
}

type ChannelReference struct {
	ID           string `json:"id" validate:"required"`
	Name         string `json:"name" readonly:"true"`
	Type         string `json:"type" readonly:"true"`
	HasRecipient bool   `json:"hasRecipient" readonly:"true"` // indicates if the channel supports/requires specifying a specific recipient
}

func (r *Rule) Validate() *errs.ErrValidation {
	err := validation.Validate.Struct(r)
	if err != nil {
		return &errs.ErrValidation{Message: err.Error()}
	}
	return nil
}
