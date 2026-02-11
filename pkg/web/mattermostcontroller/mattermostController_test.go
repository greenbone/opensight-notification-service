package mattermostcontroller

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
	mattermostChannelService := mocks.NewMattermostChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	NewMattermostController(router, notificationChannelService, mattermostChannelService, authMiddleware, registry)
	return router
}

func TestMattermostController(t *testing.T) {
	router := setup(t)

	t.Run("Create mattermost channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Post(`/notification-channel/mattermost`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Get mattermost channels is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Get(`/notification-channel/mattermost`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Update mattermost channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Putf(`/notification-channel/mattermost/%s`, uuid.New()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Delete mattermost channel is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Deletef(`/notification-channel/mattermost/%s`, uuid.New()).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})

	t.Run("Check mattermost channels is forbidden for role user", func(t *testing.T) {
		httpassert.New(t, router).
			Post(`/notification-channel/mattermost/check`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
			Expect().
			StatusCode(http.StatusForbidden)
	})
}
