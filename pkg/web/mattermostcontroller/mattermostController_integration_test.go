//go:build integration
// +build integration

package mattermostcontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MattermostController_CRUD(t *testing.T) {
	t.Parallel()

	valid := testhelper.GetValidMattermostNotificationChannel()

	t.Run("Perform all the CRUD operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		// --- Create ---
		var mattermostId string
		request.Post("/notification-channel/mattermost").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId)).
			JsonPath("$.channelName", "mattermost1").
			JsonPath("$.webhookUrl", "http://webhookurl.com/id1").
			JsonPath("$.description", "This is a test mattermost channel")
		require.NotEmpty(t, mattermostId)

		// --- List ---
		request.Get("/notification-channel/mattermost").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1)).
			JsonPath("$[0]", httpassert.HasSize(4)).
			JsonPath("$[0].id", httpassert.ExtractTo(&mattermostId)).
			JsonPath("$[0].channelName", "mattermost1").
			JsonPath("$[0].webhookUrl", "http://webhookurl.com/id1").
			JsonPath("$[0].description", "This is a test mattermost channel")

		// --- Update ---
		updated := valid
		updated.Id = &mattermostId
		newName := "updated mattermost"
		updated.ChannelName = newName
		request.Put("/notification-channel/mattermost/"+mattermostId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", mattermostId).
			JsonPath("$.channelName", newName).
			JsonPath("$.webhookUrl", "http://webhookurl.com/id1").
			JsonPath("$.description", "This is a test mattermost channel")

		// --- Delete ---
		request.Delete("/notification-channel/mattermost/" + mattermostId).
			Expect().
			StatusCode(http.StatusNoContent)

		// --- List after delete ---
		request.Get("/notification-channel/mattermost").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(0))
	})
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	mattermostSvc := notificationchannelservice.NewMattermostChannelService(svc)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewMattermostController(router, svc, mattermostSvc, testhelper.MockAuthMiddlewareWithAdmin)

	return router, db
}
