package teamsController

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamsController/teamsdto"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TeamsController_CRUD(t *testing.T) {
	t.Parallel()

	valid := testhelper.GetValidTeamsNotificationChannel()

	t.Run("Perform the Create operation", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()
		request := httpassert.New(t, router)

		var teamsId string
		request.Post("/notification-channel/teams").
			Content(`{
				"channelName": "teams1",
				"webhookUrl": "https://webhookurl.com/webhook/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			JsonPath("$.channelName", "teams1").
			JsonPath("$.webhookUrl", "https://webhookurl.com/webhook/id1").
			JsonPath("$.description", "This is a test teams channel")
		require.NotEmpty(t, teamsId)
	})

	t.Run("Perform the GET operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()
		request := httpassert.New(t, router)

		teamsId := createTeamsNotification(t, request, "teams1", valid)

		request.Get("/notification-channel/teams").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1)).
			JsonPath("$[0]", httpassert.HasSize(4)).
			JsonPath("$[0].id", httpassert.ExtractTo(&teamsId)).
			JsonPath("$[0].channelName", "teams1").
			JsonPath("$[0].webhookUrl", "https://webhookurl.com/webhook/id1").
			JsonPath("$[0].description", "This is a test teams channel")
	})

	t.Run("Perform the Update operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()
		request := httpassert.New(t, router)

		teamsId := createTeamsNotification(t, request, "teams1", valid)

		updated := valid
		newName := "updated teams"
		updated.ChannelName = newName
		request.Putf("/notification-channel/teams/%s", teamsId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(4)).
			JsonPath("$.id", teamsId).
			JsonPath("$.channelName", newName).
			JsonPath("$.webhookUrl", "https://webhookurl.com/webhook/id1").
			JsonPath("$.description", "This is a test teams channel")
	})

	t.Run("Perform the Delete operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()
		request := httpassert.New(t, router)

		teamsId := createTeamsNotification(t, request, "teams1", valid)

		request.Delete("/notification-channel/teams/" + teamsId).
			Expect().
			StatusCode(http.StatusNoContent)

		request.Get("/notification-channel/teams").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(0))
	})

	t.Run("Verify Limit check on teams limit", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 1, dummyHTTPClient())

		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
		defer func() { _ = db.Close() }()

		request := httpassert.New(t, router)

		createTeamsNotification(t, request, "teams1", valid)

		request.Post("/notification-channel/teams").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			JsonPath("$.title", "Teams channel limit reached.")
	})

	t.Run("Create two teams channels with the same name", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()
		request := httpassert.New(t, router)

		createTeamsNotification(t, request, "teams1", valid)

		request.Post("/notification-channel/teams").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusBadRequest).
			JsonPath("$.title", "Teams channel name already exists.")
	})

	t.Run("Send Test message with teams", func(t *testing.T) {
		var gotMethod, gotPath string

		rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       http.NoBody,
				Header:     make(http.Header),
			}, nil
		})

		client := http.Client{Transport: rt}

		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, client)
		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
		defer func() { _ = db.Close() }()

		request := httpassert.New(t, router)

		request.Post("/notification-channel/teams/check").
			JsonContentObject(teamsdto.TeamsNotificationChannelCheckRequest{
				WebhookUrl: "https://default41af85004428462ca2334f3ae673f7.bb.environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/01fa130f2e134641b2cf39d8a710a002/triggers/manual/paths/invoke?api-version=1&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=3WemuhZIadKcGIZdxJSmEFoJfdVJSOa1B_QNWa9z3rs",
			}).
			Expect().
			StatusCode(http.StatusNoContent)

		require.Equal(t, http.MethodPost, gotMethod)
		require.Equal(t, "/powerautomate/automations/direct/workflows/01fa130f2e134641b2cf39d8a710a002/triggers/manual/paths/invoke", gotPath)
	})
}

func createTeamsNotification(
	t *testing.T,
	request httpassert.Request,
	channelName string,
	valid teamsdto.TeamsNotificationChannelRequest,
) string {
	var teamsId string
	valid.ChannelName = channelName

	request.Post("/notification-channel/teams").
		JsonContentObject(valid).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$", httpassert.HasSize(4)).
		JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
		JsonPath("$.channelName", channelName).
		JsonPath("$.webhookUrl", "https://webhookurl.com/webhook/id1").
		JsonPath("$.description", "This is a test teams channel")
	require.NotEmpty(t, teamsId)

	return teamsId
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, dummyHTTPClient())
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)

	return router, db
}

// dummyHTTPClient returns an http.Client that always returns 204 No Content for POST requests
func dummyHTTPClient() http.Client {
	return http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       http.NoBody,
				Header:     make(http.Header),
			}, nil
		}),
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
