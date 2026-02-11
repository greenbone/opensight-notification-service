package teamscontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	notificationChannelService := mocks.NewNotificationChannelService(t)
	teamsChannelService := mocks.NewTeamsChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	NewTeamsController(router, notificationChannelService, teamsChannelService, authMiddleware, registry)
	return router
}

func TestTeamsController(t *testing.T) {
	router := setup(t)

	t.Run("Create teams channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Post(`/notification-channel/teams`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Get teams channels is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Get(`/notification-channel/teams`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Update teams channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Putf(`/notification-channel/teams/%s`, uuid.New()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Delete teams channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Deletef(`/notification-channel/teams/%s`, uuid.New()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Check teams channels is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Post(`/notification-channel/teams/check`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})
}
