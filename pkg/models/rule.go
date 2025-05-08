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

// Action determines to which channel the event is forwarded.
// Some channels (e.g. mail) require the explicit recipient(s).
type Action struct {
	Channel   ChannelReference `json:"channel" binding:"required"`
	Recipient string           `json:"recipient,omitempty"` // specific recipient if supported/required by the channel, e.g. for mail a comma separated list of mail adresses
}

type RuleOptions struct {
	Channels     []ChannelReference `json:"channels"`
	EventOrigins []OriginReference  `json:"eventOrigins"`
	EventLevels  []string           `json:"eventLevels"`
}

type OriginReference struct {
	Name      string `json:"name" readonly:"true"`
	Class     string `json:"class" binding:"required"`
	ServiceID string `json:"serviceID" binding:"required"`
}

type ChannelReference struct {
	ID           string `json:"id" binding:"required"`
	Name         string `json:"name" readonly:"true"`
	Type         string `json:"type" readonly:"true"`
	HasRecipient bool   `json:"hasRecipient" readonly:"true"` // indicates if the channel supports/requires specifying a specific recipient
}
