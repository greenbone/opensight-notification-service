// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package mattermostcontroller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupWithAuth(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	notificationChannelService := mocks.NewNotificationChannelService(t)
	mattermostChannelService := mocks.NewMattermostChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser())
	require.NoError(t, err)

	notificationChannelService.EXPECT().ListNotificationChannelsByType(mock.Anything, mock.Anything).Maybe().Return(nil, nil)
	notificationChannelService.EXPECT().DeleteNotificationChannel(mock.Anything, mock.Anything).Maybe().Return(nil, nil)

	NewMattermostController(router, notificationChannelService, mattermostChannelService, authMiddleware, registry)
	return router
}

func TestMattermostController_Permissions(t *testing.T) {
	t.Parallel()

	var endpoints = []struct {
		name   string
		method string
		path   string
	}{
		{"Create mattermost channel", http.MethodPost, "/notification-channel/mattermost"},
		{"List mattermost channels", http.MethodGet, "/notification-channel/mattermost"},
		{"Update mattermost channel", http.MethodPut, "/notification-channel/mattermost/" + uuid.NewString()},
		{"Delete mattermost channel", http.MethodDelete, "/notification-channel/mattermost/" + uuid.NewString()},
		{"Check mattermost channel", http.MethodPost, "/notification-channel/mattermost/check"},
	}

	tests := []struct {
		role      string
		wantAllow bool
	}{
		// ensure this is the same as in iam/roles.go
		{iam.OsiViewer, false},
		{iam.OsiUser, false},
		{iam.OsiAdmin, true},
		{iam.NotificationAdmin, true},
		{iam.Notification, false},
	}

	for _, tt := range tests {
		for _, ep := range endpoints {
			t.Run(ep.name+" as "+tt.role, func(t *testing.T) {
				t.Parallel()

				router := setupWithAuth(t)

				req, _ := http.NewRequest(ep.method, ep.path, nil)
				req.Header.Set("Authorization", "Bearer "+integrationTests.CreateJwtTokenWithRole(tt.role))
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if tt.wantAllow {
					require.NotEqual(t, http.StatusUnauthorized, w.Code)
					require.NotEqual(t, http.StatusForbidden, w.Code)
				} else {
					require.Equal(t, http.StatusForbidden, w.Code)
				}
			})
		}
	}
}
