// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package restErrorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/stretchr/testify/assert"
)

var someValidationError = errs.ErrValidation{
	Message: "some validation error",
	Errors:  map[string]string{"field1": "issue with field1", "field2": "issue with field2"},
}
var someConflictError = errs.ErrConflict{Errors: map[string]string{"test": "value already exists", "test2": "value already exists"}}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name              string
		err               error
		wantStatusCode    int
		wantErrorResponse *errorResponses.ErrorResponse
	}{
		{
			name:              "hide internal errors from rest clients",
			err:               errors.New("some internal error"),
			wantStatusCode:    http.StatusInternalServerError,
			wantErrorResponse: &errorResponses.ErrorInternalResponse,
		},
		{
			name:           "not found error",
			err:            fmt.Errorf("wrapped: %w", errs.ErrItemNotFound),
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:              "UnprocessableEntity error",
			err:               fmt.Errorf("wrapped: %w", &someConflictError),
			wantStatusCode:    http.StatusUnprocessableEntity,
			wantErrorResponse: helper.ToPtr(ErrConflictToResponse(someConflictError)),
		},
		{
			name:              "validation error",
			err:               fmt.Errorf("wrapped: %w", &someValidationError),
			wantStatusCode:    http.StatusBadRequest,
			wantErrorResponse: helper.ToPtr(ErrValidationToResponse(someValidationError)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(gotResponse)
			gc.Request = httptest.NewRequest(http.MethodGet, "/some/path", nil)

			ErrorHandler(gc, "some specific log message", tt.err)

			assert.Equal(t, tt.wantStatusCode, gotResponse.Code)

			if tt.wantErrorResponse != nil {
				gotErrorResponse, err := io.ReadAll(gotResponse.Body)
				if err != nil {
					t.Error("could not read response body: %w", err)
					return
				}
				wantResponseJson, err := json.Marshal(*tt.wantErrorResponse)
				if err != nil {
					t.Error("could not parse error response to json: %w", err)
				}
				assert.JSONEq(t, string(wantResponseJson), string(gotErrorResponse))
			}
		})
	}
}

func TestErrConflictToResponse(t *testing.T) {
	errConflictResponse := ErrConflictToResponse(someConflictError)

	assert.Equal(t, errorResponses.ErrorTypeGeneric, errConflictResponse.Type)
	assert.Equal(t, someConflictError.Errors, errConflictResponse.Errors)
}

func TestErrValidationToResponse(t *testing.T) {
	errValidationResponse := ErrValidationToResponse(someValidationError)

	assert.Equal(t, errorResponses.ErrorTypeValidation, errValidationResponse.Type)
	assert.Equal(t, someValidationError.Errors, errValidationResponse.Errors)
}
