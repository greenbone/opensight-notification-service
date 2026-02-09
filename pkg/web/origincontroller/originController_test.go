// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/origincontroller/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRegisterOrigins_Auth(t *testing.T) {
	tests := map[string]struct {
		authMiddleware   gin.HandlerFunc
		wantResponseCode int
	}{
		"notification user is allowed to register origins": {
			authMiddleware:   testhelper.MockAuthMiddlewareWithNotificationUser,
			wantResponseCode: http.StatusNoContent,
		},
		"normal users are not allowed to register origins": {
			authMiddleware:   testhelper.MockAuthMiddleware,
			wantResponseCode: http.StatusForbidden,
		},
		"admins are not allowed to register origins": {
			authMiddleware:   testhelper.MockAuthMiddlewareWithAdmin,
			wantResponseCode: http.StatusForbidden,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewOriginService(t)
			if tt.wantResponseCode == http.StatusNoContent {
				mockService.EXPECT().UpsertOrigins(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			}

			registry := errmap.NewRegistry()
			router := testhelper.NewTestWebEngine(registry)

			_ = NewOriginController(router, mockService, tt.authMiddleware)

			request := httpassert.New(t, router)

			request.Put("/origins/serviceX").
				JsonContentObject([]models.Origin{{Name: "Origin", Class: "origin/class"}}).
				Expect().
				StatusCode(tt.wantResponseCode)
		})
	}
}
