// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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
		wantServiceID string
		wantOrigins   []entities.Origin
		err           error
	}

	tests := map[string]struct {
		serviceID        string
		origins          []models.Origin
		serviceCalls     []ServiceCall
		wantResponseCode int
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
		},
		"invalid body (missing name)": {
			serviceID: "serviceC",
			origins: []models.Origin{
				{Name: "", Class: "origin/y"}, // Name is required
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
		},
		"invalid body (missing class)": {
			serviceID: "serviceC",
			origins: []models.Origin{
				{Name: "Origin Y", Class: ""}, // Class is required
			},
			serviceCalls:     []ServiceCall{},
			wantResponseCode: http.StatusBadRequest,
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

			_ = NewOriginController(router, mockService, testhelper.MockAuthMiddleware)

			req, err := testhelper.NewJSONRequest(http.MethodPut, "/origins/"+tt.serviceID, tt.origins)
			require.NoError(t, err, "could not build request")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantResponseCode, resp.Code)
		})
	}
}

func TestParseAndValidateOrigins(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOrigins []models.Origin
		wantErr     bool
	}{
		{
			name:  "success",
			input: `[{"name": "Origin 1", "class": "origin/1"}, {"name": "Origin 2", "class": "origin/2"}]`,
			wantOrigins: []models.Origin{
				{Name: "Origin 1", Class: "origin/1"},
				{Name: "Origin 2", Class: "origin/2"},
			},
			wantErr: false,
		},
		{
			name:    "error on empty body",
			input:   "",
			wantErr: true,
		},
		{
			name:    "error on invalid json",
			input:   `[{"name": "Origin 1"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc, _ := gin.CreateTestContext(httptest.NewRecorder())
			gc.Request = &http.Request{Body: io.NopCloser(strings.NewReader(tt.input))}

			got, err := parseAndValidateOrigins(gc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOrigins, got)
			}
		})
	}
}
