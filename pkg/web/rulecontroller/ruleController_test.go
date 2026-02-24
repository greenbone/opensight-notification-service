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
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/rulecontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *gin.Engine {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	ruleService := mocks.NewRuleService(t)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	NewRuleController(router, ruleService, authMiddleware, registry)
	return router
}

func TestRuleController(t *testing.T) {
	t.Parallel()
	setup(t)

	wantStatus := http.StatusForbidden

	tests := map[string]struct {
		method   string
		endpoint string
	}{
		"Create rule is forbidden for role user": {
			method:   http.MethodPost,
			endpoint: `/rules`,
		},
		"Get rule is forbidden for role user": {
			method:   http.MethodGet,
			endpoint: `/rules/123e4567-e89b-12d3-a456-426614174000`,
		},
		"List rules is forbidden for role user": {
			method:   http.MethodGet,
			endpoint: `/rules`,
		},
		"Update rule is forbidden for role user": {
			method:   http.MethodPut,
			endpoint: `/rules/123e4567-e89b-12d3-a456-426614174000`,
		},
		"Delete rule is forbidden for role user": {
			method:   http.MethodDelete,
			endpoint: `/rules/123e4567-e89b-12d3-a456-426614174000`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			httpassert.New(t, setup(t)).
				Perform(tt.method, tt.endpoint).
				AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
				Expect().
				StatusCode(wantStatus)
		})
	}
}
