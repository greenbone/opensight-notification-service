// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/require"
)

func createRule(t *testing.T, router *gin.Engine, ruleJSON string) (ruleID string) {
	httpassert.New(t, router).Post("/rules").
		AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
		JsonContent(ruleJSON).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$.id", httpassert.ExtractTo(&ruleID))
	require.NotEmpty(t, ruleID)
	return ruleID
}

func Test_UpdateRule(t *testing.T) {
	t.Parallel()

	t.Run("update all values", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{
			{Name: "origin0", Class: "serviceA/origin0"},
			{Name: "origin1", Class: "serviceA/origin1"},
		}
		channels := []models.NotificationChannel{
			{ChannelName: new("channel-name-0"), ChannelType: "mattermost"},
			{ChannelName: new("channel-name-1"), ChannelType: "mail"},
		}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": { "id": "%s" }
				}
			}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Updated Test Rule",
				"trigger": {
					"levels": ["urgent"],
					"origins": [{
						"name": "read-only,ignored",
						"class": "serviceA/origin1",
						"serviceID": "read-only,ignored"
					}]
				},
				"action": {
					"channel": {
						"id": "%s",
						"name": "read-only,ignored",
						"type": "read-only,ignored"
					},
					"recipient": "test@example.org"
				},
				"active": true
			}`, *channels[1].Id)).
			Expect().
			StatusCode(http.StatusOK).
			JsonTemplate(`{
				"id": "<value>",
				"name": "Updated Test Rule",
				"trigger": {
					"levels": ["urgent"],
					"origins": [{
						"name": "origin1",
						"class": "serviceA/origin1",
						"serviceID": "<value>"
					}]
				},
				"action": {
					"channel": {
						"id": "<value>",
						"name": "channel-name-1",
						"type": "mail"
					},
					"recipient": "test@example.org"
				},
				"active": true
			}`,
				map[string]any{
					"$.id":                           ruleID,
					"$.trigger.origins[0].serviceID": origins[1].ServiceID,
					"$.action.channel.id":            *channels[1].Id,
				},
			)
	})

	t.Run("failure due to missing required fields", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": { "id": "%s" }
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{}`). // missing required fields
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"name": "A name is required.",
					"trigger.origins": "At least one origin is required.",
					"trigger.levels": "At least one level is required.",
					"action.channel.id": "A channel is required."
				}
			}`)
	})

	t.Run("failure due to invalid field values", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": { "id": "%s" }
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"trigger": {
					"levels": [""],
					"origins": [{}]
				},
				"action": {
					"channel": {"id": "invalid-uuid"}
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
					"action.channel.id": "Channel ID must be a valid UUIDv4."
				}
			}`)
	})

	t.Run("failure due to missing recipient for mail channel", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("mail-1"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"},
					"recipient": "a@example.com"
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule Updated",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"},
					"recipient": ""
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
		channels := []models.NotificationChannel{{ChannelName: new("mattermost-1"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule Updated",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"},
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
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Updated Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "non-existent" }]
				},
				"action": {
					"channel": {"id": "%s"}
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

	t.Run("failure due to non-existing channel id", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": { "id": "%s" }
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"name": "Updated Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "839bbc13-24ec-4079-be42-63af9a9b66ac"}
				}
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"action.channel.id": "Channel does not exist."
				}
			}`)
	})

	// already existing rule name
	t.Run("failure due to already existing rule name", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		createRule(t, router, fmt.Sprintf(`{
				"name": "Other Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
		}`, *channels[0].Id))

		ruleID := createRule(t, router, fmt.Sprintf(`{
				"name": "Test Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
		}`, *channels[0].Id))

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Other Rule",
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
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

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		httpassert.New(t, router).Putf("/rules/%s", uuid.New().String()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(fmt.Sprintf(`{
				"name": "Test Rule",
				
				"trigger": {
					"levels": ["info"],
					"origins": [{ "class": "serviceA/origin0" }]
				},
				"action": {
					"channel": {"id": "%s"}
				}
			}`, *channels[0].Id)).
			Expect().
			StatusCode(http.StatusNotFound)
	})

	t.Run("failure due to invalid id", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		updateRule := models.Rule{
			Name: "Test Rule",

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

		httpassert.New(t, router).Put("/rules/invalid-id").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updateRule).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type":"greenbone/validation-error",
			  	"title":"ID must be a valid UUIDv4."
			}`)
	})
}
