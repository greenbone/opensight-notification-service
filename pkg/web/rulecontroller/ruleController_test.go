// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulecontroller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/rulecontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupWithAuth(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	ruleService := mocks.NewRuleService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser())
	require.NoError(t, err)

	ruleService.EXPECT().List(mock.Anything).Maybe().Return([]models.Rule{}, nil)
	ruleService.EXPECT().GetAllRuleOptions(mock.Anything).Maybe().Return(&models.RuleOptions{}, nil)
	ruleService.EXPECT().Get(mock.Anything, mock.Anything).Maybe().Return(models.Rule{}, nil)
	ruleService.EXPECT().Delete(mock.Anything, mock.Anything).Maybe().Return(nil)

	NewRuleController(router, ruleService, authMiddleware, registry)
	return router
}

func TestRuleController_Permissions(t *testing.T) {
	t.Parallel()

	var endpoints = []struct {
		name   string
		method string
		path   string
	}{
		{"Create rule", http.MethodPost, "/rules"},
		{"Get rule", http.MethodGet, "/rules/123e4567-e89b-12d3-a456-426614174000"},
		{"List rules", http.MethodGet, "/rules"},
		{"Update rule", http.MethodPut, "/rules/123e4567-e89b-12d3-a456-426614174000"},
		{"Delete rule", http.MethodDelete, "/rules/123e4567-e89b-12d3-a456-426614174000"},
		{"Rule options", http.MethodGet, "/rules/ruleoptions"},
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
