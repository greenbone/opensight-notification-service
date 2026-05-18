// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulecontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/rulecontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*gin.Engine, *mocks.RuleService) {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	ruleService := mocks.NewRuleService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	ruleService.On("List", mock.Anything).Maybe().Return([]models.Rule{}, nil)
	ruleService.On("GetAllRuleOptions", mock.Anything).Maybe().Return(&models.RuleOptions{}, nil)

	NewRuleController(router, ruleService, authMiddleware, registry)
	return router, ruleService
}

func TestRuleController_ForbiddenRoles(t *testing.T) {
	t.Parallel()

	forbiddenRoles := []string{iam.OsiViewer, iam.User, iam.OsiUser, iam.OsiAdmin, iam.Notification}

	endpoints := []struct {
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

	for _, role := range forbiddenRoles {
		for _, ep := range endpoints {
			t.Run(ep.name+" is forbidden for role "+role, func(t *testing.T) {
				t.Parallel()

				router, _ := setup(t)

				httpassert.New(t, router).
					Perform(ep.method, ep.path).
					AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
					Expect().
					StatusCode(http.StatusForbidden)
			})
		}
	}
}

// We use a group for the whole ruleController, testing only the one requiring no request body.
func TestRuleController_AllowedRoles(t *testing.T) {
	t.Parallel()

	allowedRoles := []string{iam.Admin, iam.NotificationAdmin}

	for _, role := range allowedRoles {
		t.Run("List rules is allowed for role "+role, func(t *testing.T) {
			t.Parallel()

			router, _ := setup(t)

			httpassert.New(t, router).
				Get("/rules").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Expect().
				StatusCode(http.StatusOK)
		})

		t.Run("Rule options is allowed for role "+role, func(t *testing.T) {
			t.Parallel()

			router, _ := setup(t)

			httpassert.New(t, router).
				Get("/rules/ruleoptions").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Expect().
				StatusCode(http.StatusOK)
		})
	}
}
