// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/origincontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterOrigins(t *testing.T) {
	type ServiceCall struct {
		wantServiceID string
		wantOrigins   []entities.Origin
		err           error
	}

	tests := map[string]struct {
		serviceID        string
		origins          []models.Origin
		serviceCalls     []ServiceCall
		wantResponseCode int
		wantBodyContains string
	}{
		"valid request": {
			serviceID: "serviceA",
			origins: []models.Origin{
				{Name: "Origin 1", Class: "origin/1"},
				{Name: "Origin 2", Class: "origin/2"},
			},
			serviceCalls: []ServiceCall{
				{
					wantServiceID: "serviceA",
					wantOrigins: []entities.Origin{
						{Name: "Origin 1", Class: "origin/1"},
						{Name: "Origin 2", Class: "origin/2"},
					},
					err: nil,
				},
			},
			wantResponseCode: http.StatusNoContent,
		},
		"service error": {
			serviceID: "serviceB",
			origins: []models.Origin{
				{Name: "Origin X", Class: "origin/x"},
			},
			serviceCalls: []ServiceCall{
				{
					wantServiceID: "serviceB",
					wantOrigins: []entities.Origin{
						{Name: "Origin X", Class: "origin/x"},
					},
					err: errors.New("internal error"),
				},
			},
			wantResponseCode: http.StatusInternalServerError,
			wantBodyContains: "internal",
		},
		"invalid body (missing name)": {
			serviceID: "serviceC",
			origins: []models.Origin{
				{Name: "", Class: "origin/y"},
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
			wantBodyContains: "Name",
		},
		"invalid body (missing class)": {
			serviceID: "serviceC",
			origins: []models.Origin{
				{Name: "Origin Y", Class: ""},
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
			wantBodyContains: "Class",
		},
		"empty origins list": {
			serviceID: "serviceD",
			origins:   []models.Origin{},
			serviceCalls: []ServiceCall{
				{
					wantServiceID: "serviceD",
					wantOrigins:   []entities.Origin{},
					err:           nil,
				},
			},
			wantResponseCode: http.StatusNoContent,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewOriginService(t)

			for _, call := range tt.serviceCalls {
				mockService.EXPECT().UpsertOrigins(mock.Anything, call.wantServiceID, call.wantOrigins).Return(call.err).Once()
			}

			registry := errmap.NewRegistry()
			router := testhelper.NewTestWebEngine(registry)

			_ = NewOriginController(router, mockService, testhelper.MockAuthMiddlewareWithNotificationUser)

			request := httpassert.New(t, router)

			resp := request.Put("/origins/" + tt.serviceID).
				JsonContentObject(tt.origins).
				Expect().
				StatusCode(tt.wantResponseCode)
			if tt.wantBodyContains != "" {
				assert.Contains(t, resp.GetBody(), tt.wantBodyContains)
			}
		})
	}
}

func setupWithAuth(t *testing.T) *gin.Engine {
	mockService := mocks.NewOriginService(t)
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	NewOriginController(router, mockService, authMiddleware)
	return router
}

func TestRegisterOrigins_Permissions(t *testing.T) {
	t.Parallel()

	var endpoints = []struct {
		name   string
		method string
		path   string
	}{
		{"Register origins", http.MethodPut, "/origins/serviceX"},
	}

	tests := []struct {
		role      string
		wantAllow bool
	}{
		// ensure this is the same as in iam/roles.go
		{iam.OsiViewer, false},
		{iam.User, false},
		{iam.OsiUser, false},
		{iam.OsiAdmin, false},
		{iam.Admin, false},
		{iam.NotificationAdmin, false},
		{iam.Notification, true},
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
