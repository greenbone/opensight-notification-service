// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/origincontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterOrigins(t *testing.T) {
	type ServiceCall struct {
		wantNamespace string
		wantOrigins   []entities.Origin
		err           error
	}

	tests := map[string]struct {
		namespace        string
		origins          []models.Origin
		serviceCalls     []ServiceCall
		wantResponseCode int
	}{
		"valid request": {
			namespace: "serviceA",
			origins: []models.Origin{
				{Name: "Origin 1", Class: "origin/1"},
				{Name: "Origin 2", Class: "origin/2"},
			},
			serviceCalls: []ServiceCall{
				{
					wantNamespace: "serviceA",
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
			namespace: "serviceB",
			origins: []models.Origin{
				{Name: "Origin X", Class: "origin/x"},
			},
			serviceCalls: []ServiceCall{
				{
					wantNamespace: "serviceB",
					wantOrigins: []entities.Origin{
						{Name: "Origin X", Class: "origin/x"},
					},
					err: errors.New("internal error"),
				},
			},
			wantResponseCode: http.StatusInternalServerError,
		},
		"invalid body (missing name)": {
			namespace: "serviceC",
			origins: []models.Origin{
				{Name: "", Class: "origin/y"}, // Name is required
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
		},
		"invalid body (missing class)": {
			namespace: "serviceC",
			origins: []models.Origin{
				{Name: "Origin Y", Class: ""}, // Class is required
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
		},
		"empty origins list": {
			namespace: "serviceD",
			origins:   []models.Origin{},
			serviceCalls: []ServiceCall{
				{
					wantNamespace: "serviceD",
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
				mockService.EXPECT().UpsertOrigins(mock.Anything, call.wantNamespace, call.wantOrigins).Return(call.err).Once()
			}

			registry := errmap.NewRegistry()
			router := testhelper.NewTestWebEngine(registry)

			_ = NewOriginController(router, mockService, testhelper.MockAuthMiddleware)

			req, err := testhelper.NewJSONRequest(http.MethodPut, "/origins/"+tt.namespace, tt.origins)
			require.NoError(t, err, "could not build request")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantResponseCode, resp.Code)
		})
	}
}
