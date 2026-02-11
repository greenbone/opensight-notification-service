package usesCases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/stretchr/testify/require"
)

func TestListTeamsChannels(t *testing.T) {
	t.Run("List teams channels", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		var teamsId string

		// Create teams channel
		request.Post("/notification-channel/teams").
			JsonContent(`{
				"channelName": "teams1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			JsonTemplate(`{
				"id": "d9cc9be2-7b4d-4c6f-991d-a40cfe002ceb",
				"channelName": "teams1",
				"webhookUrl": "https://example.com/hooks/id1",
				"description": "This is a test teams channel"
			}`, map[string]any{
				"id": httpassert.IgnoreJsonValue,
			})
		require.NotEmpty(t, teamsId)

		// List teams channels
		request.Get("/notification-channel/teams").
			Expect().
			StatusCode(http.StatusOK).
			JsonTemplate(`[
				{
					"id": "fb46613b-4bf8-45c7-ad6f-e83e5ced8b81",
					"channelName": "teams1",
					"webhookUrl": "https://example.com/hooks/id1",
					"description": "This is a test teams channel"
				}
			]`, map[string]any{
				"$.0.id": httpassert.IgnoreJsonValue,
			})
	})
}
