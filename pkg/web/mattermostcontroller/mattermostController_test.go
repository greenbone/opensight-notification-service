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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	notificationChannelService := mocks.NewNotificationChannelService(t)
	mattermostChannelService := mocks.NewMattermostChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	// We only test permissions, the method itself is not part of these tests.
	notificationChannelService.
		On("ListNotificationChannelsByType", mock.Anything, mock.Anything).
		Maybe().
		Return(nil, nil)

	NewMattermostController(router, notificationChannelService, mattermostChannelService, authMiddleware, registry)
	return router
}

func TestMattermostController_ForbiddenRoles(t *testing.T) {
	router := setup(t)

	forbiddenRoles := []string{iam.OsiViewer, iam.User, iam.OsiUser, iam.OsiAdmin, iam.Notification}

	type endpoint struct {
		name   string
		method string
		path   string
	}

	endpoints := []endpoint{
		{"Create mattermost channel", http.MethodPost, "/notification-channel/mattermost"},
		{"List mattermost channels", http.MethodGet, "/notification-channel/mattermost"},
		{"Update mattermost channel", http.MethodPut, "/notification-channel/mattermost/" + uuid.NewString()},
		{"Delete mattermost channel", http.MethodDelete, "/notification-channel/mattermost/" + uuid.NewString()},
		{"Check mattermost channel", http.MethodPost, "/notification-channel/mattermost/check"},
	}

	for _, role := range forbiddenRoles {
		for _, ep := range endpoints {
			t.Run(ep.name+" is forbidden for role "+role, func(t *testing.T) {
				httpassert.New(t, router).
					Perform(ep.method, ep.path).
					AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
					Expect().
					StatusCode(http.StatusForbidden)
			})
		}
	}
}

func TestMattermostController_AllowedRoles(t *testing.T) {
	router := setup(t)

	allowedRoles := []string{iam.Admin, iam.NotificationAdmin}

	for _, role := range allowedRoles {
		t.Run("Access is granted for role "+role, func(t *testing.T) {
			httpassert.New(t, router).
				Get(`/notification-channel/mattermost`).
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Expect().
				StatusCode(http.StatusOK)
		})
	}
}
