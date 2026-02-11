package usesCases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMattermostChannel(t *testing.T) {
	t.Run("Update mattermost channel", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		var mattermostId string

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId))
		require.NotEmpty(t, mattermostId)

		// Update mattermost channel
		request.Putf("/notification-channel/mattermost/%s", mattermostId).
			JsonContent(`{
				"channelName": "mattermost2",
				"webhookUrl": "https://example.com/hooks/id2",
				"description": "This is a test mattermost channel changed"
			}`).
			Expect().
			StatusCode(http.StatusOK).
			JsonTemplate(`{
				"id": "fb46613b-4bf8-45c7-ad6f-e83e5ced8b81",
				"channelName": "mattermost2",
				"webhookUrl": "https://example.com/hooks/id2",
				"description": "This is a test mattermost channel changed"
			}`, map[string]any{
				"$.id": mattermostId,
			})
	})

	t.Run("Update mattermost channel with an invalid webhook URL returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		var mattermostId string

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId))
		require.NotEmpty(t, mattermostId)

		// Update mattermost channel
		request.Putf("/notification-channel/mattermost/%s", mattermostId).
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

	t.Run("Update mattermost channel without required fields returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		var mattermostId string

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId))
		require.NotEmpty(t, mattermostId)

		// Update mattermost channel
		request.Putf("/notification-channel/mattermost/%s", mattermostId).
			JsonContent(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Log().
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"channelName": "A channel name is required.",
					"webhookUrl": "A Webhook URL is required."
				}
			}`)
	})

	t.Run("Update mattermost channel name with an existing one returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		var mattermostId string

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated)

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost 2",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId))
		require.NotEmpty(t, mattermostId)

		// Update mattermost channel
		request.Putf("/notification-channel/mattermost/%s", mattermostId).
			JsonContent(`{
				"channelName": "mattermost 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Log().
			Json(`{
				"type": "greenbone/generic-error",
				"title": "Mattermost channel name already exists."
			}`)
	})
}
