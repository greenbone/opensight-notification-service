// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRule(t *testing.T, router *gin.Engine, rule models.Rule) (ruleID string) {
	httpassert.New(t, router).Post("/rules").
		AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
		JsonContentObject(rule).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$.id", httpassert.ExtractTo(&ruleID))
	require.NotEmpty(t, ruleID)
	return ruleID
}

func Test_UpdateRule(t *testing.T) {
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

		originalRule := models.Rule{
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

		updatedRule := models.Rule{
			Name: "Updated Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"urgent"},
				Origins: []models.OriginReference{{Name: "read-only,ignored", Class: "serviceA/origin1", ServiceID: "read-only,ignored"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: *channels[1].Id,
				},
				Recipient: "test@example.com",
			},
			Active: true,
		}
		wantUpdatedRule := models.Rule{
			Name: "Updated Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"urgent"},
				Origins: []models.OriginReference{{Name: "origin1", Class: "serviceA/origin1", ServiceID: origins[1].ServiceID}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID:   *channels[1].Id,
					Name: "channel-name-1",
					Type: "mail",
				},
				Recipient: "test@example.com",
			},
			Active: true,
		}

		ruleID := createRule(t, router, originalRule)

		resp := httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
			Expect().
			StatusCode(http.StatusOK)
		var gotUpdatedRule models.Rule
		resp.GetJsonBodyObject(&gotUpdatedRule)

		assert.NotEmpty(t, gotUpdatedRule.ID)
		gotUpdatedRule.ID = ""
		assert.Equal(t, wantUpdatedRule, gotUpdatedRule)
	})

	t.Run("failure due to missing required fields", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		originalRule := models.Rule{
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
		updatedRule := models.Rule{} // missing required fields

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		originalRule := models.Rule{
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

		updatedRule := models.Rule{
			Trigger: models.Trigger{
				Origins: []models.OriginReference{{}},
				Levels:  []string{""},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: "invalid-uuid",
				},
			},
		}

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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
		channels := []models.NotificationChannel{{ChannelName: new("mail-1"), ChannelType: "mail"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		originalRule := models.Rule{
			Name: "Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel:   models.ChannelReference{ID: *channels[0].Id},
				Recipient: "a@example.com",
			},
		}

		updatedRule := models.Rule{
			Name: "Test Rule Updated",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel:   models.ChannelReference{ID: *channels[0].Id},
				Recipient: "", // missing but required for mail channel
			},
		}

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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

		originalRule := models.Rule{
			Name: "Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{ID: *channels[0].Id},
			},
		}

		updatedRule := models.Rule{
			Name: "Test Rule Updated",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel:   models.ChannelReference{ID: *channels[0].Id},
				Recipient: "not@supported.com",
			},
		}

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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

		originalRule := models.Rule{
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
		updatedRule := models.Rule{
			Name: "Updated Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "non-existent"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: *channels[0].Id,
				},
			},
		}

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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

		originalRule := models.Rule{
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
		updatedRule := models.Rule{
			Name: "Updated Test Rule",
			Trigger: models.Trigger{
				Levels:  []string{"info"},
				Origins: []models.OriginReference{{Class: "serviceA/origin0"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{
					ID: uuid.NewString(), // non-existing channel ID
				},
			},
		}

		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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

	// already existing rule name
	t.Run("failure due to already existing rule name", func(t *testing.T) {
		t.Parallel()
		ruleLimit := 10
		origins := []entities.Origin{{Name: "origin0", Class: "serviceA/origin0"}}
		channels := []models.NotificationChannel{{ChannelName: new("channel-name-0"), ChannelType: "mattermost"}}
		router := setupTestEnvironment(t, origins, channels, ruleLimit)

		otherUntouchedRule := models.Rule{
			Name: "Other Rule",
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
		originalRule := models.Rule{
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
		updatedRule := models.Rule{
			Name: "Other Rule", // name already used by other rule

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

		createRule(t, router, otherUntouchedRule)
		ruleID := createRule(t, router, originalRule)

		httpassert.New(t, router).Putf("/rules/%s", ruleID).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updatedRule).
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

		httpassert.New(t, router).Putf("/rules/%s", uuid.New().String()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContentObject(updateRule).
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
