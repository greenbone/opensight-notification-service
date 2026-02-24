// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
)

func Test_ListRules(t *testing.T) {
	t.Parallel()

	t.Run("get empty list of rules", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Get("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK).
			Json("[]")
	})

	t.Run("get all rules ordered by name", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		createRule(t, router, fmt.Sprintf(`{
				"name": "Rule B",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
			}`, *channels[0].Id))

		createRule(t, router, fmt.Sprintf(`{
				"name": "Rule A",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
			}`, *channels[0].Id))

		httpassert.New(t, router).Get("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$[0].id", httpassert.NotEmpty()).
			JsonPath("$[1].id", httpassert.NotEmpty()).
			JsonTemplate(`[
					{
						"id": "<value>",
						"name": "Rule A",
						"trigger": {
							"levels": ["info"],
							"origins": [
								{ 
									"name": "origin0", 
									"class": "serviceA/origin0", 
									"serviceID": "<value>" 
								}
							]
						},
						"action": {
							"channel": {
								"id": "<value>",
								"name": "channel-name",
								"type": "mattermost"
							}
						},
						"active": false
					},
					{
						"id": "<value>",
						"name": "Rule B",
						"trigger": {
							"levels": ["info"],
							"origins": [
								{ 
									"name": "origin0", 
									"class": "serviceA/origin0", 
									"serviceID": "<value>" 
								}
							]
						},
						"action": {
							"channel": {
								"id": "<value>",
								"name": "channel-name",
								"type": "mattermost"
							}
						},
						"active": false
					}
				]`,
				map[string]any{
					"$.0.id":                           httpassert.IgnoreJsonValue,
					"$.0.trigger.origins[0].serviceID": origins[0].ServiceID,
					"$.0.action.channel.id":            *channels[0].Id,
					"$.1.id":                           httpassert.IgnoreJsonValue,
					"$.1.trigger.origins[0].serviceID": origins[0].ServiceID,
					"$.1.action.channel.id":            *channels[0].Id,
				},
			)
	})
}
