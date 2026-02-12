package usesCases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/require"
)

func TestUpdateTeamsChannel(t *testing.T) {
	t.Run("Update teams channel", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		var teamsId string

		// Create teams channel
		httpassert.New(t, router).Post("/notification-channel/teams").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId))
		require.NotEmpty(t, teamsId)

		// Update teams channel
		httpassert.New(t, router).Putf("/notification-channel/teams/%s", teamsId).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams2",
				"webhookUrl": "https://example.com/hooks/id2",
				"description": "This is a test teams channel changed"
			}`).
			Expect().
			StatusCode(http.StatusOK).
			JsonTemplate(`{
				"id": "fb46613b-4bf8-45c7-ad6f-e83e5ced8b81",
				"channelName": "teams2",
				"webhookUrl": "https://example.com/hooks/id2",
				"description": "This is a test teams channel changed"
			}`, map[string]any{
				"$.id": teamsId,
			})
	})

	t.Run("Update teams channel with an invalid webhook URL returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		var teamsId string

		// Create teams channel
		httpassert.New(t, router).Post("/notification-channel/teams").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId))
		require.NotEmpty(t, teamsId)

		// Update teams channel
		httpassert.New(t, router).Putf("/notification-channel/teams/%s", teamsId).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "1",
				"webhookUrl": "invalid",
				"description": "b"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"webhookUrl": "Please enter a valid webhook URL."
				}
			}`)
	})

	t.Run("Update teams channel without required fields returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		var teamsId string

		// Create teams channel
		httpassert.New(t, router).Post("/notification-channel/teams").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId))
		require.NotEmpty(t, teamsId)

		// Update teams channel
		httpassert.New(t, router).Putf("/notification-channel/teams/%s", teamsId).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"channelName": "A channel name is required.",
					"webhookUrl": "A Webhook URL is required."
				}
			}`)
	})

	t.Run("Update teams channel name with an existing one returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		var teamsId string

		// Create teams channel
		httpassert.New(t, router).Post("/notification-channel/teams").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated)

		// Create teams channel
		httpassert.New(t, router).Post("/notification-channel/teams").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams 2",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId))
		require.NotEmpty(t, teamsId)

		// Update teams channel
		httpassert.New(t, router).Putf("/notification-channel/teams/%s", teamsId).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"channelName": "teams 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "MS Teams channel name already exists."
			}`)
	})
}
