// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/require"
)

func Test_GetRule(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [
						{ 
							"name": "read-only,ignored", 
							"class": "serviceA/origin0", 
							"serviceID": "read-only-ignored" 
						}
					]
				},
				"action": {
					"channel": {
						"id": "%s",
						"name": "read-only,ignored",
						"type": "read-only,ignored"
					},
					"recipient": "a@example.com"
				},
				"active": true
		}`, *channels[0].Id))

		httpassert.New(t, router).Getf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.id", httpassert.NotEmpty()).
			JsonTemplate(`{
				"id": "<value>",
				"name": "Test Rule",
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
						"name": "channel-name-0",
						"type": "mail"
					},
					"recipient": "a@example.com"
				},
				"active": true
			}`,
				map[string]any{
					"$.id":                           httpassert.IgnoreJsonValue,
					"$.trigger.origins[0].serviceID": origins[0].ServiceID,
					"$.action.channel.id":            *channels[0].Id,
				})
	})

	t.Run("get invalid rule as deactived and with errors set", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mail"}}
		router, channelRepo, originRepo := setupTestEnvironmentWithRepoReturn(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"},
					"recipient": "a@example.com"
				},
				"active": true
		}`, *channels[0].Id))

		// delete the origin and channel to make the rule invalid
		ctx := context.Background()
		err := originRepo.UpsertOrigins(ctx, origins[0].ServiceID, []entities.Origin{})
		require.NoError(t, err)
		err = channelRepo.DeleteNotificationChannel(ctx, *channels[0].Id)
		require.NoError(t, err)

		httpassert.New(t, router).Getf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.id", httpassert.IsUUID()).
			JsonTemplate(`{
				"id": "<value>",
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": []
				},
				"action": {
					"channel": {
						"id": "",
						"name": "",
						"type": ""
					},
					"recipient": "a@example.com"
				},
				"active": false,
				"errors": {
					"trigger.origins": "At least one origin is required.",
					"action.channel.id": "A channel is required."
				}
			}`,
				map[string]any{
					"$.id": httpassert.IgnoreJsonValue,
				})
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Getf("/rules/%s", uuid.NewString()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNotFound)
	})

	t.Run("failure due to invalid id", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Get("/rules/invalid-uuid").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type":"greenbone/validation-error",
			  	"title":"ID must be a valid UUIDv4."
			}`)
	})
}
