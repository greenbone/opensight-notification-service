// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package restErrorHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-golang-libraries/pkg/logs"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
)

// ErrorHandler determines the appropriate error response and code from the error type. It relies on the types defined in [errs].
// The default case is an internal server error hiding the implementation details from the client. In this case a log message is issued containing the error.
// A log message for context can be provided via parameter internalErrorLogMessage.
func ErrorHandler(gc *gin.Context, internalErrorLogMessage string, err error) {
	var errConflict *errs.ErrConflict
	var errValidation *errs.ErrValidation
	switch {
	case errors.Is(err, errs.ErrItemNotFound):
		gc.JSON(http.StatusNotFound, errorResponses.NewErrorGenericResponse(err.Error()))
	case errors.As(err, &errConflict):
		gc.JSON(http.StatusUnprocessableEntity, ErrConflictToResponse(*errConflict))
	case errors.As(err, &errValidation):
		gc.JSON(http.StatusBadRequest, ErrValidationToResponse(*errValidation))
	default:
		logs.Ctx(gc.Request.Context()).Err(err).Str("endpoint", gc.Request.Method+" "+gc.Request.URL.Path).Msg(internalErrorLogMessage)
		gc.JSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
	}
}

func ErrValidationToResponse(err errs.ErrValidation) errorResponses.ErrorResponse {
	return errorResponses.NewErrorValidationResponse(err.Message, "", err.Errors)
}

func ErrConflictToResponse(err errs.ErrConflict) errorResponses.ErrorResponse {
	return errorResponses.ErrorResponse{
		Type:   errorResponses.ErrorTypeGeneric,
		Title:  err.Message,
		Errors: err.Errors,
	}
}
