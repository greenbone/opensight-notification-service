// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package mailcontroller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*gin.Engine, *mocks.MailChannelService) {
	registry := errmap.NewRegistry()
	engine := testhelper.NewTestWebEngine(registry)

	notificationChannelServicer := mocks.NewMailChannelService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	AddCheckMailServerController(engine, notificationChannelServicer, authMiddleware, registry)
	return engine, notificationChannelServicer
}

func TestCheckMailServer(t *testing.T) {
	t.Run("mail server check is successful", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(nil)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": true,
				"isTlsEnforced": false,
				"username": "testUser",
				"password": "123"
			}`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.NotificationAdmin)).
			Expect().
			StatusCode(http.StatusNoContent).
			NoContent()
	})

	t.Run("none request body", func(t *testing.T) {
		engine, _ := setup(t)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`-`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.NotificationAdmin)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "unable to parse the request",
				"details":"error parsing body"
			}`)
	})

	t.Run("minimal required fields", func(t *testing.T) {
		engine, _ := setup(t)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{}`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.NotificationAdmin)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"domain": "A mailhub is required.",
					"port": "A port is required."
				}
			}`)
	})

	t.Run("username and password are required if isAuthenticationRequired is true", func(t *testing.T) {
		engine, _ := setup(t)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": true,
				"isTlsEnforced": false,
				"username": "",
				"password": ""
			}`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.NotificationAdmin)).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"username": "A username is required.",
					"password": "A password is required."
				}
			}`)
	})

	t.Run("return the error if the mail server check fails", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(notificationchannelservice.ErrMailServerUnreachable)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": false,
				"isTlsEnforced": false
			}`).
			AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.NotificationAdmin)).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "Server is unreachable"
			}`)
	})
}

func TestCheckMailServer_Permissions(t *testing.T) {
	t.Parallel()

	var endpoints = []struct {
		name   string
		method string
		path   string
	}{
		{"Create mail channel", http.MethodPost, "/notification-channel/mail/check"},
	}

	tests := []struct {
		role      string
		wantAllow bool
	}{
		// ensure this is the same as in iam/roles.go
		{iam.OsiViewer, false},
		{iam.User, false},
		{iam.OsiUser, false},
		{iam.OsiAdmin, true},
		{iam.Admin, true},
		{iam.NotificationAdmin, true},
		{iam.Notification, false},
	}

	for _, tt := range tests {
		for _, ep := range endpoints {
			t.Run(ep.name+" as "+tt.role, func(t *testing.T) {
				t.Parallel()

				router, notificationChannelServicer := setup(t)
				notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).Maybe().Return(nil)

				req, _ := http.NewRequest(ep.method, ep.path, strings.NewReader(`{"domain":"example.com","port":123}`))
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
