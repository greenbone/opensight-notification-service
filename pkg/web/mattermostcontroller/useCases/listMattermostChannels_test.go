package usesCases

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/stretchr/testify/require"
)

func TestListMattermostChannels(t *testing.T) {
	t.Run("List mattermost channels", func(t *testing.T) {
		t.Parallel()

		router, db := setupTestRouter(t)
		defer db.Close()

		var mattermostId string

		// Create mattermost channel
		httpassert.New(t, router).Post("/notification-channel/mattermost").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
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

		// List mattermost channels
		httpassert.New(t, router).Get("/notification-channel/mattermost").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusOK).
			JsonTemplate(`[
				{
					"id": "fb46613b-4bf8-45c7-ad6f-e83e5ced8b81",
					"channelName": "mattermost1",
					"webhookUrl": "https://example.com/hooks/id1",
					"description": "This is a test mattermost channel"
				}
			]`, map[string]any{
				"$.0.id": httpassert.IgnoreJsonValue,
			})
	})
}
