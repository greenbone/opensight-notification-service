package teamsController

import (
	"net/http"
	"testing"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TeamsController_CRUD(t *testing.T) {
	t.Parallel()

	const TeamsGetQuery = "SELECT id, channel_name, webhook_url, description FROM notification_service.notification_channel WHERE id = $1"

	t.Run("Perform the Create operation", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, http.Client{Timeout: 15 * time.Second})
		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
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
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			StatusCode(http.StatusCreated)
		require.NotEmpty(t, teamsId)

		var dbEntry dbTeamsChannel
		err := db.Get(&dbEntry, TeamsGetQuery, teamsId)
		require.NoError(t, err, "DB entry for created teams channel should exist")
		require.Equal(t, "teams1", dbEntry.ChannelName)
		require.Equal(t, "https://webhookurl.com/webhook/id1", dbEntry.WebhookUrl)
		require.Equal(t, "This is a test teams channel", dbEntry.Description)
	})

	t.Run("Perform the Update operations", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, http.Client{Timeout: 15 * time.Second})
		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
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
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			StatusCode(http.StatusCreated)
		require.NotEmpty(t, teamsId)

		request.Putf("/notification-channel/teams/%s", teamsId).
			Content(`{
				"channelName": "updated teams",
				"webhookUrl": "https://webhookurl.com/webhook/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusOK)

		var dbEntry dbTeamsChannel
		err := db.Get(&dbEntry, TeamsGetQuery, teamsId)
		require.NoError(t, err, "DB entry for created teams channel should exist")
		require.Equal(t, "updated teams", dbEntry.ChannelName)
		require.Equal(t, "https://webhookurl.com/webhook/id1", dbEntry.WebhookUrl)
		require.Equal(t, "This is a test teams channel", dbEntry.Description)

	})

	t.Run("Perform the Delete operations", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, http.Client{Timeout: 15 * time.Second})
		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
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
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			StatusCode(http.StatusCreated)

		request.Delete("/notification-channel/teams/" + teamsId).
			Expect().
			StatusCode(http.StatusNoContent)

		var dbEntry dbTeamsChannel
		err := db.Get(&dbEntry, TeamsGetQuery, teamsId)
		require.Error(t, err, "sql: no rows in result set")
	})

	t.Run("Verify Limit check on teams limit", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 1, http.Client{Timeout: 15 * time.Second})

		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
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
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			StatusCode(http.StatusCreated)

		request.Post("/notification-channel/teams").
			Content(`{
				"channelName": "teams1",
				"webhookUrl": "https://webhookurl.com/webhook/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			JsonPath("$.title", "Teams channel limit reached.")
	})

	t.Run("Create two teams channels with the same name", func(t *testing.T) {
		repo, db := testhelper.SetupNotificationChannelTestEnv(t)
		svc := notificationchannelservice.NewNotificationChannelService(repo)
		teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, http.Client{Timeout: 15 * time.Second})
		registry := errmap.NewRegistry()
		router := testhelper.NewTestWebEngine(registry)
		AddTeamsController(router, svc, teamsSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)
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
			JsonPath("$.id", httpassert.ExtractTo(&teamsId)).
			StatusCode(http.StatusCreated)

		request.Post("/notification-channel/teams").
			Content(`{
				"channelName": "teams1",
				"webhookUrl": "https://webhookurl.com/webhook/id1",
				"description": "This is a test teams channel"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			JsonPath("$.title", "Teams channel name already exists.")
	})
}

type dbTeamsChannel struct {
	ID          string `db:"id"`
	ChannelName string `db:"channel_name"`
	WebhookUrl  string `db:"webhook_url"`
	Description string `db:"description"`
}
