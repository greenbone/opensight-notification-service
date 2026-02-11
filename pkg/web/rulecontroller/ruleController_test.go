// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulecontroller

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/originrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/rulerepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/ruleservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T) (
	*gin.Engine,
	notificationrepository.NotificationChannelRepository,
	*originrepository.OriginRepository,
) {
	notificationChannelRepo, db := testhelper.SetupNotificationChannelTestEnv(t)

	originRepo, err := originrepository.NewOriginRepository(db)
	require.NoError(t, err)

	ruleRepo, err := rulerepository.NewRuleRepository(db)
	require.NoError(t, err)
	ruleService := ruleservice.NewRuleService(ruleRepo, 10)

	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	_ = NewRuleController(router, ruleService, testhelper.MockAuthMiddlewareWithAdmin, registry)

	return router, notificationChannelRepo, originRepo
}

func TestIntegration_Create_Notification(t *testing.T) {

	serviceID := "serviceA"
	origin := entities.Origin{
		Name:  "tada",
		Class: "origin/tada",
	}

	tests := map[string]struct {
		ruleIn           models.Rule
		wantStatusCode   int
		wantBodyContains string
	}{
		"failure due to non existing notification origin": {
			ruleIn: models.Rule{
				Name: "Test Rule",

				Trigger: models.Trigger{
					Origins: []models.OriginReference{models.OriginReference{Class: "non-existing"}},
				},
			},
			wantStatusCode:   400,
			wantBodyContains: "origin1",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			router, notificationChannelRepo, originRepo := setupTestRouter(t)

			err := originRepo.UpsertOrigins(ctx, serviceID, []entities.Origin{origin})
			require.NoError(t, err)

			channel, err := notificationChannelRepo.CreateNotificationChannel(ctx, models.NotificationChannel{
				ChannelName: helper.ToPtr("mattermost-test"),
				ChannelType: "mattermost",
				WebhookUrl:  helper.ToPtr("url"),
			})
			require.NoError(t, err)
			channelId := channel.Id
			_ = channelId

			req := httpassert.New(t, router)

			resp := req.Post("/rules").
				JsonContentObject(tt.ruleIn).
				Expect().
				StatusCode(tt.wantStatusCode)

			if tt.wantBodyContains != "" {
				assert.Contains(t, resp.GetBody(), tt.wantBodyContains)
			}
		})
	}
}
