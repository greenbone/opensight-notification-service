// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package restErrorHandler

import "github.com/gin-gonic/gin"

// TODO: 26.01.2026 stolksdorf - eduardo's branch will fix this
// ErrorHandler determines the appropriate error response and code from the error type. It relies on the types defined in [errs].
// The default case is an internal server error hiding the implementation details from the client. In this case a log message is issued containing the error.
// A log message for context can be provided via parameter internalErrorLogMessage.
func ErrorHandler(gc *gin.Context, internalErrorLogMessage string, err error) {
	//	var errConflict *errs.ErrConflict
	//	var errValidation *errs.ErrValidation
	//	switch {
	//	case errors.Is(err, errs.ErrItemNotFound):
	//		gc.JSON(http.StatusNotFound, errorResponses.NewErrorGenericResponse(err.Error()))
	//	case errors.As(err, &errConflict):
	//		gc.JSON(http.StatusUnprocessableEntity, ErrConflictToResponse(*errConflict))
	//	case errors.As(err, &errValidation):
	//		gc.JSON(http.StatusBadRequest, ErrValidationToResponse(*errValidation))
	//	default:
	//		logs.Ctx(gc.Request.Context()).Err(err).Str("endpoint", gc.Request.Method+" "+gc.Request.URL.Path).Msg(internalErrorLogMessage)
	//		gc.JSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
	//	}
}

//	func ErrValidationToResponse(err errs.ErrValidation) errorResponses.ErrorResponse {
//		return errorResponses.NewErrorValidationResponse(err.Message, "", err.Errors)
//	}
//
//	func ErrConflictToResponse(err errs.ErrConflict) errorResponses.ErrorResponse {
//		return errorResponses.ErrorResponse{
//			Type:   errorResponses.ErrorTypeGeneric,
//			Title:  err.Message,
//			Errors: err.Errors,
//		}
//	}
func NotificationChannelErrorHandler(gc *gin.Context, title string, errs map[string]string, err error) {
	//	if len(errs) > 0 && title != "" {
	//		gc.JSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse(title, "", errs))
	//		return
	//	}
	//
	//	switch {
	//	case errors.Is(err, mailcontroller.ErrMailChannelBadRequest) ||
	//		errors.Is(err, mattermostcontroller.ErrMattermostChannelBadRequest) ||
	//		errors.Is(err, teamsController.ErrTeamsChannelBadRequest):
	//		gc.JSON(http.StatusBadRequest,
	//			errorResponses.NewErrorValidationResponse("Invalid mail channel data.", "", nil))
	//	case errors.Is(err, notificationchannelservice.ErrListMailChannels) ||
	//		errors.Is(err, notificationchannelservice.ErrListMattermostChannels) ||
	//		errors.Is(err, notificationchannelservice.ErrListTeamsChannels):
	//		gc.JSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
	//	case errors.Is(err, notificationchannelservice.ErrMattermostChannelNameExists) ||
	//		errors.Is(err, notificationchannelservice.ErrTeamsChannelNameExists):
	//		gc.JSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("Channel name should be unique.", "",
	//			map[string]string{"channelName": "Channel name should be unique."}))
	//	case errors.Is(err, notificationchannelservice.ErrMailChannelLimitReached):
	//		gc.JSON(http.StatusUnprocessableEntity, errorResponses.NewErrorValidationResponse("Mail channel limit reached.", "",
	//			map[string]string{"channelName": "Mail channel already exists."}))
	//	case errors.Is(err, notificationchannelservice.ErrMattermostChannelLimitReached):
	//		gc.JSON(http.StatusUnprocessableEntity, errorResponses.NewErrorValidationResponse("Mattermost channel limit reached.", "",
	//			map[string]string{"channelName": "Mattermost channel creation limit reached."}))
	//	case errors.Is(err, notificationchannelservice.ErrTeamsChannelLimitReached):
	//		gc.JSON(http.StatusUnprocessableEntity, errorResponses.NewErrorValidationResponse("Teams channel limit reached.", "",
	//			map[string]string{"channelName": "Teams channel creation limit reached."}))
	//	default:
	//		gc.JSON(http.StatusInternalServerError, errorResponses.ErrorInternalResponse)
	//	}
}
