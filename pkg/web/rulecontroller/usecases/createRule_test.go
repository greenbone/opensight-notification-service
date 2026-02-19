// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/originrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/rulerepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/ruleservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/rulecontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment creates the given origins and notification channels and returns a router
// to be used for the rule endpoint tests. No mocks are used in this setup.
// The passed `origins` and `channels` slices are populated with the read-only fields.
// IMPORTANT: If you run tests in parallel, you must not pass the same instance of a slice in multiple tests
// as they are modified by this function.
func setupTestEnvironment(t *testing.T, origins []entities.Origin, channels []models.NotificationChannel, ruleLimit int) *gin.Engine {
	router, _, _ := setupTestEnvironmentWithRepoReturn(t, origins, channels, ruleLimit)
	return router
}

func setupTestEnvironmentWithRepoReturn(t *testing.T, origins []entities.Origin, channels []models.NotificationChannel, ruleLimit int) (
	*gin.Engine,
	notificationrepository.NotificationChannelRepository,
	*originrepository.OriginRepository,
) {
	ctx := context.Background()

	// create notification channels
	notificationChannelRepo, db := testhelper.SetupNotificationChannelTestEnv(t)
	for i := range channels {
		channel, err := notificationChannelRepo.CreateNotificationChannel(ctx, channels[i])
		require.NoError(t, err)
		require.NotNil(t, channel.Id)
		channels[i].Id = channel.Id
	}

	// create origins
	originRepo, err := originrepository.NewOriginRepository(db)
	require.NoError(t, err)

	serviceID := "serviceA"
	err = originRepo.UpsertOrigins(ctx, serviceID, origins)
	require.NoError(t, err)
	for i := range origins {
		origins[i].ServiceID = serviceID
	}

	ruleRepo, err := rulerepository.NewRuleRepository(db)
	require.NoError(t, err)
	ruleService := ruleservice.NewRuleService(ruleRepo, notificationChannelRepo, ruleLimit)

	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	rulecontroller.NewRuleController(router, ruleService, authMiddleware, registry)

	return router, notificationChannelRepo, originRepo
}

func Test_CreateRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{
						"name": "read-only-ignored",
						"class": "serviceA/origin0",
						"serviceID": "read-only-ignored"
					}]
				},
				"action": {
					"channel": {
						"id": "%s",
						"name": "read-only-ignored",
						"type": "read-only-ignored"
					},
					"recipient": "a@example.com"
				},
				"Active": true
	}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.NotEmpty()).
			JsonTemplate(`{
				"id": "<value>",
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{
						"name": "origin0",
						"class": "serviceA/origin0",
						"serviceID": "<value>"
					}]
				},
				"action": {
					"channel": {
						"id": "<value>",
						"name": "channel-name",
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
				},
			)
	})

	t.Run("failure due to missing required fields", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"name": "A name is required.",
					"trigger.origins": "At least one origin is required.",
					"trigger.levels": "At least one level is required.",
					"trigger.action.channel.id": "A channel is required."
				}
			}`)
	})

	t.Run("failure due to invalid field values", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		router := setupTestEnvironment(t, nil, nil, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"trigger": {
					"origins": [{}],
					"levels": [""]
				},
				"action": {
					"channel": {
						"id": "invalid-uuid"
					}
				}
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"name": "A name is required.",
					"trigger.origins[0].class": "An origin class is required.",
					"trigger.levels[0]": "A level is required.",
					"trigger.action.channel.id": "Channel ID must be a valid UUIDv4."
				}
			}`)
	})

	t.Run("failure due to missing recipient for mail channel", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
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
						
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"trigger.action.recipient": "Recipient is required for the selected channel."
				}
			}`)
	})

	t.Run("failure due to non-empty recipient for non-mail channel", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {
						"id": "%s"
					},
					"recipient": "not@supported.com"
				}
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"trigger.action.recipient": "Recipient is not supported for the selected channel."
				}
			}`)
	})

	t.Run("failure due to non existing notification origin", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, nil, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "non-existing" }]
				},
				"action": {
					"channel": {
						"id": "%s"
					}
				}
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"trigger.origins": "One or more origins do not exist."
				}
			}`)
	})

	t.Run("failure due to non existing channel id", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		router := setupTestEnvironment(t, origins, nil, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {
						"id": "9e9912cf-97be-491e-8d7e-93a992334a3a"
					}
				}
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"trigger.action.channel.id": "Channel does not exist."
				}
			}`)
	})

	t.Run("failure due to already existing rule name", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
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
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusCreated)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
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
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"name": "Alert rule name already exists."
				}
			}`)
	})

	t.Run("failure due to limit of alert rules reached", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 1
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
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
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusCreated)

		httpassert.New(t, router).Post("/rules").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule 2",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {
						"id": "%s"
					}						
				}
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "Alert rule limit reached."
			}`)
	})
}
