//go:build integration
// +build integration

package mattermostcontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/request"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MattermostController_CRUD(t *testing.T) {
	t.Parallel()

	valid := testhelper.GetValidMattermostNotificationChannel()

	t.Run("Perform the Create operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		var mattermostId string
		request.Post("/notification-channel/mattermost").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", httpassert.ExtractTo(&mattermostId)).
			JsonPath("$.channelName", "mattermost1").
			JsonPath("$.webhookUrl", "https://webhookurl.com/hooks/id1").
			JsonPath("$.description", "This is a test mattermost channel")
		require.NotEmpty(t, mattermostId)
	})

	t.Run("Perform the GET operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		mattermostId := createMattermostNotification(t, request, "mattermost1", valid)

		request.Get("/notification-channel/mattermost").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1)).
			JsonPath("$[0]", httpassert.HasSize(4)).
			JsonPath("$[0].id", httpassert.ExtractTo(&mattermostId)).
			JsonPath("$[0].channelName", "mattermost1").
			JsonPath("$[0].webhookUrl", "https://webhookurl.com/hooks/id1").
			JsonPath("$[0].description", "This is a test mattermost channel")
	})

	t.Run("Perform the Update operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		mattermostId := createMattermostNotification(t, request, "mattermost1", valid)

		updated := valid
		newName := "updated mattermost"
		updated.ChannelName = newName
		request.Putf("/notification-channel/mattermost/%s", mattermostId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", mattermostId).
			JsonPath("$.channelName", newName).
			JsonPath("$.webhookUrl", "https://webhookurl.com/hooks/id1").
			JsonPath("$.description", "This is a test mattermost channel")
	})

	t.Run("Perform the Delete operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		mattermostId := createMattermostNotification(t, request, "mattermost1", valid)

		request.Delete("/notification-channel/mattermost/" + mattermostId).
			Expect().
			StatusCode(http.StatusNoContent)

		request.Get("/notification-channel/mattermost").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(0))
	})

	t.Run("Verify Limit check on mattermost limit", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		mattermostSvc := notificationchannelservice.NewMattermostChannelService(svc, 1)
		gin.SetMode(gin.TestMode)
		router := gin.New()
		NewMattermostController(router, svc, mattermostSvc, testhelper.MockAuthMiddlewareWithAdmin)
		defer db.Close()

		request := httpassert.New(t, router)

		createMattermostNotification(t, request, "mattermost1", valid)

		request.Post("/notification-channel/mattermost").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			JsonPath("$.title", "Mattermost channel limit reached.")
	})

	t.Run("Create two mattermost channels with the same name", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)

		createMattermostNotification(t, request, "mattermost1", valid)

		request.Post("/notification-channel/mattermost").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusBadRequest).
			JsonPath("$.title", "Channel name should be unique.")
	})
}

func createMattermostNotification(t *testing.T, request httpassert.Request, channelName string, valid request.MattermostNotificationChannelRequest) string {
	var mattermostId string
	valid.ChannelName = channelName

	request.Post("/notification-channel/mattermost").
		JsonContentObject(valid).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$", httpassert.HasSize(4)).
		JsonPath("$.id", httpassert.ExtractTo(&mattermostId)).
		JsonPath("$.channelName", channelName).
		JsonPath("$.webhookUrl", "https://webhookurl.com/hooks/id1").
		JsonPath("$.description", "This is a test mattermost channel")
	require.NotEmpty(t, mattermostId)

	return mattermostId
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	mattermostSvc := notificationchannelservice.NewMattermostChannelService(svc, 20)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewMattermostController(router, svc, mattermostSvc, testhelper.MockAuthMiddlewareWithAdmin)

	return router, db
}
