package usesCases

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
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

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	teamscontroller.NewTeamsController(router, svc, teamsSvc, authMiddleware, registry)

	return router, db
}

func TestDeleteTeamsChannel(t *testing.T) {
	t.Run("Delete a teams channel without proper role returns error", func(t *testing.T) {
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
		httpassert.New(t, router).Deletef("/notification-channel/teams/%s", teamsId).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			Expect().
			StatusCode(http.StatusNoContent)
	})
}
