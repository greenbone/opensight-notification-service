package usesCases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/stretchr/testify/require"
)

func TestCreateMattermostChannel(t *testing.T) {
	t.Run("Create mattermost channel", func(t *testing.T) {
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
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId)).
			JsonTemplate(`{
				"id": "d9cc9be2-7b4d-4c6f-991d-a40cfe002ceb",
				"channelName": "mattermost1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`, map[string]any{
				"id": httpassert.IgnoreJsonValue,
			})
		require.NotEmpty(t, mattermostId)
	})

	t.Run("Create mattermost channel with invalid webhook URL returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "a",
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

	t.Run("Create mattermost channel without required fields returns an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
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

	t.Run("Create mattermost channel with an existing name return an error", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		// Create mattermost channel
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated)

		// Create mattermost channel with the same name
		request.Post("/notification-channel/mattermost").
			JsonContent(`{
				"channelName": "mattermost 1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test mattermost channel"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "Mattermost channel name already exists."
			}`)
	})
}
