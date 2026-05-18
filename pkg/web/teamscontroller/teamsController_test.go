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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	notificationChannelService := mocks.NewNotificationChannelService(t)
	teamsChannelService := mocks.NewTeamsChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	notificationChannelService.
		On("ListNotificationChannelsByType", mock.Anything, mock.Anything).
		Maybe().
		Return(nil, nil)

	NewTeamsController(router, notificationChannelService, teamsChannelService, authMiddleware, registry)
	return router
}

func TestTeamsController_ForbiddenRoles(t *testing.T) {
	router := setup(t)

	forbiddenRoles := []string{iam.OsiViewer, iam.User, iam.OsiUser, iam.OsiAdmin, iam.Notification}

	endpoints := []struct {
		name   string
		method string
		path   string
	}{
		{"Create teams channel", http.MethodPost, "/notification-channel/teams"},
		{"List teams channels", http.MethodGet, "/notification-channel/teams"},
		{"Update teams channel", http.MethodPut, "/notification-channel/teams/" + uuid.NewString()},
		{"Delete teams channel", http.MethodDelete, "/notification-channel/teams/" + uuid.NewString()},
		{"Check teams channel", http.MethodPost, "/notification-channel/teams/check"},
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

func TestTeamsController_AllowedRoles(t *testing.T) {
	router := setup(t)

	allowedRoles := []string{iam.Admin, iam.NotificationAdmin}

	for _, role := range allowedRoles {
		t.Run("Access is granted for role "+role, func(t *testing.T) {
			httpassert.New(t, router).
				Get(`/notification-channel/teams`).
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Expect().
				StatusCode(http.StatusOK)
		})
	}
}
