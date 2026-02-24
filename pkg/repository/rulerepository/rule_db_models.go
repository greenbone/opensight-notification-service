// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulerepository

import (
	"bytes"
	"encoding/json"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/lib/pq"
)

const ruleTable = "notification_service.rules"
const channelTable = "notification_service.notification_channel"
const originsTable = "notification_service.origins"

func ruleSelectWithJoin(from string) string {
	return `SELECT
		r.id,
		r.name,
		r.trigger_origins,
		r.trigger_levels,
		r.action_channel_id,
		r.action_recipient,
		r.active,
		c.channel_name,
		c.channel_type,
		COALESCE(
			json_agg(
				json_build_object(
					'name', o.name,
					'class', o.class,
					'serviceID', o.service_id
				)
			) FILTER (WHERE o.class IS NOT NULL),
			CAST('[]' AS json)
		) AS origins
	FROM ` + from + ` r
	LEFT JOIN ` + channelTable + ` c ON r.action_channel_id = c.id
	LEFT JOIN ` + originsTable + ` o ON o.class = ANY(r.trigger_origins)
	`
}

var ruleQuerySelect = ruleSelectWithJoin(ruleTable)

const ruleQueryGroupBy = `
GROUP BY r.id, r.name, r.trigger_origins, r.trigger_levels, r.action_channel_id, r.action_recipient, r.active, c.channel_name, c.channel_type`

var createRuleQuery = `WITH inserted AS (
		INSERT INTO ` + ruleTable + ` (
			name, trigger_origins, trigger_levels, action_channel_id, action_recipient, active
		) VALUES (
			:name, :trigger_origins, :trigger_levels, :action_channel_id, :action_recipient, :active
		)
		RETURNING id, name, trigger_origins, trigger_levels, action_channel_id, action_recipient, active
	)
` + ruleSelectWithJoin("inserted") + ruleQueryGroupBy

var getRuleByIdQuery = ruleQuerySelect + ` WHERE r.id = $1 ` + ruleQueryGroupBy

var listRulesUnfilteredQuery = ruleQuerySelect + ruleQueryGroupBy + ` ORDER BY r.name`

var updateRuleQuery = `WITH updated AS (
UPDATE ` + ruleTable + `
	SET name = :name,
		trigger_origins = :trigger_origins,
		trigger_levels = :trigger_levels,
		action_channel_id = :action_channel_id,
		action_recipient = :action_recipient,
		active = :active
	WHERE id = :id
	RETURNING id, name, trigger_origins, trigger_levels, action_channel_id, action_recipient, active
)
` + ruleSelectWithJoin("updated") + ruleQueryGroupBy

const deleteQuery = `DELETE FROM ` + ruleTable + ` WHERE id = $1`

type ruleRow struct {
	ID              string         `db:"id"`
	Name            string         `db:"name"`
	TriggerOrigins  pq.StringArray `db:"trigger_origins"`
	TriggerLevels   pq.StringArray `db:"trigger_levels"`
	ActionChannelID string         `db:"action_channel_id"`
	ActionRecipient *string        `db:"action_recipient"`
	Active          bool           `db:"active"`
	channelRow
	originRow
}

// columns joined from notification_channel table
type channelRow struct {
	ChannelName string `db:"channel_name"`
	ChannelType string `db:"channel_type"`
}

// data joined from origins table
type originRow struct {
	OriginsJSON json.RawMessage `db:"origins"`
}

// ToModel converts a ruleRow to a models.Rule
func (r ruleRow) ToModel() (models.Rule, error) {
	// Unmarshal origins JSON
	var originsParsed []models.OriginReference
	if len(r.OriginsJSON) > 0 && string(r.OriginsJSON) != "null" {
		jsonDecoder := json.NewDecoder(bytes.NewReader(r.OriginsJSON))
		jsonDecoder.DisallowUnknownFields() // to ensure the correct fields are used in the db query json_agg
		if err := jsonDecoder.Decode(&originsParsed); err != nil {
			return models.Rule{}, err
		}
	}

	rule := models.Rule{
		ID:   r.ID,
		Name: r.Name,
		Trigger: models.Trigger{
			Origins: originsParsed,
			Levels:  []string(r.TriggerLevels),
		},
		Action: models.Action{
			Channel: models.ChannelReference{
				ID:   r.ActionChannelID,
				Name: r.ChannelName,
				Type: r.ChannelType,
			},
			Recipient: helper.SafeDereference(r.ActionRecipient),
		},
		Active: r.Active,
	}

	return rule, nil
}

// toRuleRow converts a models.Rule to a ruleRow for insert
func toRuleRow(rule models.Rule) ruleRow {
	// Extract origin classes (the only writable field)
	originClasses := make([]string, 0, len(rule.Trigger.Origins))
	for _, origin := range rule.Trigger.Origins {
		originClasses = append(originClasses, origin.Class)
	}

	row := ruleRow{
		Name:            rule.Name,
		TriggerOrigins:  originClasses,
		TriggerLevels:   rule.Trigger.Levels,
		ActionChannelID: rule.Action.Channel.ID, // take only the writable field
		ActionRecipient: helper.ToPtr(rule.Action.Recipient),
		Active:          rule.Active,
	}

	return row
}
