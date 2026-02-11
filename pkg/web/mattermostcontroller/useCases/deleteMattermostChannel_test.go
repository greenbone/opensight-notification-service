package usesCases

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	mattermostSvc := notificationchannelservice.NewMattermostChannelService(svc, 20, &http.Client{Timeout: 15 * time.Second})
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	mattermostcontroller.NewMattermostController(router, svc, mattermostSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)

	return router, db
}

func TestDeleteMattermostChannel(t *testing.T) {
	t.Run("Delete a mattermost channel without proper role returns error", func(t *testing.T) {
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
			Log().
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

		// Delete mattermost channel
		request.Deletef("/notification-channel/mattermost/%s", mattermostId).
			Expect().
			StatusCode(http.StatusNoContent)
	})
}
