// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

// A rule determines which events cause which action.
// Each incoming event is matched with the trigger conditions.
// If the condition is fulfilled, the provided action is triggered.
type Rule struct {
	ID      string  `json:"id" readonly:"true"`
	Name    string  `json:"name" binding:"required"`
	Trigger Trigger `json:"trigger" binding:"required"`
	Action  Action  `json:"action" binding:"required"`
	Active  bool    `json:"active"`
}

// Trigger condition, fulfilled if both one of `origins` and `levels` match the ones from the incomming event.
type Trigger struct {
	Origins []OriginReference `json:"origins" binding:"required"`
	Levels  []string          `json:"level" binding:"required"`
}

// Action determines to which sink the event is forwarded.
// Some sinks (e.g. mail) require the explicit recipient(s).
type Action struct {
	Sink      SinkReference `json:"sink" binding:"required"`
	Recipient string        `json:"recipient,omitempty"` // specific recipient if supported/required by the sink, e.g. for mail a comma separated list of mail adresses
}

type RuleOptions struct {
	Sinks        []SinkReference   `json:"sinks"`
	EventOrigins []OriginReference `json:"eventOrigins"`
	EventLevels  []string          `json:"eventLevels"`
}
