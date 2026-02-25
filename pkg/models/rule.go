// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import (
	"fmt"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/translation"
	"github.com/greenbone/opensight-notification-service/pkg/validation"
)

// A rule determines which events cause which action.
// Each incoming event is matched with the trigger conditions.
// If the condition is fulfilled, the provided action is triggered.
type Rule struct {
	ID      string           `json:"id" readonly:"true"`
	Name    string           `json:"name" validate:"required"`
	Trigger Trigger          `json:"trigger" validate:"required"`
	Action  Action           `json:"action" validate:"required"`
	Active  bool             `json:"active"`
	Errors  ValidationErrors `json:"errors,omitempty" readonly:"true"` // populated if the rule is invalid, this can be useful to highlight rules which need action from the user.
}

// Trigger condition, fulfilled if both one of `origins` and `levels` match the ones from the incoming event.
type Trigger struct {
	Origins []OriginReference `json:"origins" validate:"required"`
	Levels  []string          `json:"levels" validate:"required"`
}

// Action determines to which channel the event is forwarded.
// Some channels (e.g. mail) require the explicit recipient(s).
type Action struct {
	Channel   ChannelReference `json:"channel" validate:"required"`
	Recipient string           `json:"recipient,omitempty"` // specific recipient if supported/required by the channel, e.g. for mail a comma separated list of mail adresses
}

type OriginReference struct {
	Name      string `json:"name" readonly:"true"`
	Class     string `json:"class" validate:"required"`
	ServiceID string `json:"serviceID" readonly:"true"`
}

type ChannelReference struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" readonly:"true"`
	Type string `json:"type" readonly:"true"`
}

func (r *Rule) Cleanup() {
	r.Name = strings.TrimSpace(r.Name)
}

// Validate checks if the rule is valid and returns validation errors if not.
// Furthermore it populates the `Errors` field with these validation errors.
func (r *Rule) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if r.Name == "" {
		errs["name"] = translation.NameIsRequired
	}

	if len(r.Trigger.Origins) == 0 {
		errs["trigger.origins"] = translation.OriginsAreRequired
	} else {
		for i, origin := range r.Trigger.Origins {
			if origin.Class == "" {
				errs[fmt.Sprintf("trigger.origins[%d].class", i)] = translation.OriginClassIsRequired
			}
		}
	}

	if len(r.Trigger.Levels) == 0 {
		errs["trigger.levels"] = translation.LevelsAreRequired
	} else {
		for i, level := range r.Trigger.Levels {
			if level == "" {
				errs[fmt.Sprintf("trigger.levels[%d]", i)] = translation.LevelIsRequired
			}
		}
	}

	if r.Action.Channel.ID == "" {
		errs["action.channel.id"] = translation.ChannelIsRequired
	} else {
		err := validation.Validate.Var(r.Action.Channel.ID, "uuid4")
		if err != nil {
			errs["action.channel.id"] = translation.InvalidChannelID
		}
	}

	if len(errs) > 0 {
		r.Errors = errs
		return errs
	}

	return nil
}
