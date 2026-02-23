// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
)

func Test_DeleteRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{
			{Name: "origin0", Class: "serviceA/origin0"},
		}
		channels := []models.NotificationChannel{
			{ChannelName: new("channel-name-0"), ChannelType: "mattermost"},
		}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {
						"id": "%s"
					}
				}
			}`, *channels[0].Id))

		httpassert.New(t, router).Deletef("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNoContent)

		// verify rule is deleted
		httpassert.New(t, router).Getf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNotFound)
	})

	t.Run("deleting non-existing rule is a no-op", func(t *testing.T) {
		t.Parallel()
		router := setupTestEnvironment(t, nil, nil, 10)

		httpassert.New(t, router).Deletef("/rules/%s", uuid.NewString()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNoContent)
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
