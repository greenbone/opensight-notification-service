package usesCases

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamscontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, &http.Client{Timeout: 15 * time.Second})
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	teamscontroller.NewTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)

	return router, db
}

func TestDeleteTeamsChannel(t *testing.T) {
	t.Run("Delete a teams channel without proper role returns error", func(t *testing.T) {
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

		// Delete teams channel
		request.Deletef("/notification-channel/teams/%s", teamsId).
			Expect().
			StatusCode(http.StatusNoContent)
	})
}
