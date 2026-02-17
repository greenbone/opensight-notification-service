// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
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
	t.Run("success", func(t *testing.T) {
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		rule := models.Rule{
			Name: "Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Name: "read-only,ignored", Class: "serviceA/origin0", ServiceID: "read-only-ignored"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID:   *channels[0].Id,
					Name: "read-only,ignored",
					Type: "read-only,ignored",
				},
				Recipient: "a@example.com",
			},
		}
		wantRule := models.Rule{
			Name: "Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Name: "origin0", Class: "serviceA/origin0", ServiceID: origins[0].ServiceID}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID:   *channels[0].Id,
					Name: "channel-name-0",
					Type: "mail",
				},
				Recipient: "a@example.com",
			},
		}

		ruleID := createRule(t, router, rule)
		wantRule.ID = ruleID

		var gotRule models.Rule
		resp := httpassert.New(t, router).Getf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK)

		resp.GetJsonBodyObject(&gotRule)
		require.Equal(t, wantRule, gotRule)
	})

	t.Run("not found", func(t *testing.T) {
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Getf("/rules/%s", uuid.NewString()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNotFound)
	})

	t.Run("failure due to invalid id", func(t *testing.T) {
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
