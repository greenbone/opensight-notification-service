package usesCases

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamscontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func setup(t *testing.T, transport http.Client) *gin.Engine {
	t.Helper()

	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	teamsSvc := notificationchannelservice.NewTeamsChannelService(svc, 20, &transport)
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	teamscontroller.NewTeamsController(router, svc, teamsSvc, authMiddleware, registry)
	defer db.Close()
	return router
}

func TestCheckTeamsChannel(t *testing.T) {
	t.Run("Check teams channel with new webhook format", func(t *testing.T) {
		t.Parallel()

		transport := http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
					Header:     make(http.Header),
				}, nil
			}),
		}
		router := setup(t, transport)

		// Check teams channel
		httpassert.New(t, router).Post("/notification-channel/teams/check").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"webhookUrl": "https://example.com/hooks/id1"
			}`).
			Expect().
			StatusCode(http.StatusNoContent)
	})

	t.Run("Check teams channel with old webhook format", func(t *testing.T) {
		t.Parallel()

		transport := http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
					Header:     make(http.Header),
				}, nil
			}),
		}
		router := setup(t, transport)

		// Check teams channel
		httpassert.New(t, router).Post("/notification-channel/teams/check").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"webhookUrl": "https://example.com/webhook/id1"
			}`).
			Expect().
			StatusCode(http.StatusNoContent)
	})

	t.Run("Check teams channel with invalid webhook URL returns an error", func(t *testing.T) {
		t.Parallel()

		transport := http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
					Header:     make(http.Header),
				}, nil
			}),
		}
		router := setup(t, transport)

		// Check teams channel
		httpassert.New(t, router).Post("/notification-channel/teams/check").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"webhookUrl": "invalid"
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

	t.Run("Check teams channel with teams server response 404", func(t *testing.T) {
		t.Parallel()

		transport := http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       http.NoBody,
					Header:     make(http.Header),
					Status:     "404 Not Found",
				}, nil
			}),
		}
		router := setup(t, transport)

		// Check teams channel
		httpassert.New(t, router).Post("/notification-channel/teams/check").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{
				"webhookUrl": "https://example.com/hooks/id1"
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "teams message could not be send: http status: 404 Not Found"
			}`)
	})

	t.Run("Check teams channel without required url", func(t *testing.T) {
		t.Parallel()

		transport := http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       http.NoBody,
					Header:     make(http.Header),
					Status:     "404 Not Found",
				}, nil
			}),
		}
		router := setup(t, transport)

		// Check teams channel
		httpassert.New(t, router).Post("/notification-channel/teams/check").
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
			JsonContent(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"title":"",
				"type":"greenbone/validation-error",
				"errors": {
					"webhookUrl":"A Webhook URL is required."
				}
			}`)
	})
}
