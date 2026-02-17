// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/assert"
)

func Test_ListRules(t *testing.T) {
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

		rule1 := models.Rule{
			Name: "Rule A",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: *channels[0].Id,
				},
			},
		}
		rule2 := models.Rule{
			Name: "Rule B",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: *channels[0].Id,
				},
			},
		}
		wantRule1 := models.Rule{
			Name: "Rule A",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Name: "origin0", Class: "serviceA/origin0", ServiceID: origins[0].ServiceID}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID:   *channels[0].Id,
					Name: "channel-name",
					Type: "mattermost",
				},
			},
		}
		wantRule2 := models.Rule{
			Name: "Rule B",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Name: "origin0", Class: "serviceA/origin0", ServiceID: origins[0].ServiceID}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID:   *channels[0].Id,
					Name: "channel-name",
					Type: "mattermost",
				},
			},
		}

		ruleID2 := createRule(t, router, rule2)
		ruleID1 := createRule(t, router, rule1)
		// set ids for comparison
		wantRule1.ID = ruleID1
		wantRule2.ID = ruleID2

		wantRules := []models.Rule{wantRule1, wantRule2} // alphabetical order by name

		resp := httpassert.New(t, router).Get("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK)

		var gotRules []models.Rule
		resp.GetJsonBodyObject(&gotRules)

		assert.Equal(t, wantRules, gotRules)
	})
}
