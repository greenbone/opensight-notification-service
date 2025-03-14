// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"errors"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
)

func getNotification() models.Notification {
	return models.Notification{
		Id:        "57fe22b8-89a4-445f-b6c7-ef9ea724ea48",
		Timestamp: time.Time{}.Format(time.RFC3339Nano),
		Origin:    "Example Task XY",
		Title:     "Example Task XY failed",
		Detail:    "Example Task XY failed because ...",
		Level:     "error",
	}
}

func TestListNotifications(t *testing.T) {
	someNotification := getNotification()

	type mockReturn struct {
		items        []models.Notification
		totalResults uint64
		err          error
	}
	type want struct {
		serviceCall    bool
		serviceArg     query.ResultSelector
		responseCode   int
		responseParsed query.ResponseListWithMetadata[models.Notification]
	}

	resultSelectorWithoutFilter := query.ResultSelector{
		Filter: &filter.Request{Operator: filter.LogicOperatorAnd},
		Paging: &paging.Request{PageSize: 10},
		Sorting: &sorting.Request{
			SortColumn:    dtos.OccurrenceFieldName,
			SortDirection: sorting.DirectionAscending,
		},
	}

	tests := []struct {
		name        string
		requestBody query.ResultSelector
		mockReturn  mockReturn
		want        want
	}{
		{
			name:        "service is called with correct result selector",
			requestBody: resultSelectorWithoutFilter,
			mockReturn: mockReturn{
				items:        []models.Notification{someNotification},
				totalResults: 1,
				err:          nil,
			},
			want: want{
				serviceCall:  true,
				serviceArg:   resultSelectorWithoutFilter,
				responseCode: http.StatusOK,
				responseParsed: query.ResponseListWithMetadata[models.Notification]{
					Data:     []models.Notification{someNotification},
					Metadata: query.NewMetadata(resultSelectorWithoutFilter, 1),
				},
			},
		},
		{
			name:        "return internal server error on service failure",
			requestBody: resultSelectorWithoutFilter,
			mockReturn:  mockReturn{err: errors.New("internal service error")},
			want: want{
				serviceCall:  true,
				serviceArg:   resultSelectorWithoutFilter,
				responseCode: http.StatusInternalServerError,
			},
		},
		{
			name:        "return bad request on invalid input",
			requestBody: query.ResultSelector{Paging: &paging.Request{PageSize: -1}}, // invalid page size
			want: want{
				serviceCall:  false,
				responseCode: http.StatusBadRequest,
			},
		},
	}

	requestUrl := "/notifications"

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotificationService := mocks.NewNotificationService(t)

			// Create a new engine for testing
			engine := gin.Default()
			// constructor registers the routes
			_ = NewNotificationController(&engine.RouterGroup, mockNotificationService, testhelper.MockAuthMiddleware)

			if tt.want.serviceCall {
				mockNotificationService.EXPECT().ListNotifications(mock.Anything, tt.want.serviceArg).
					Return(tt.mockReturn.items, tt.mockReturn.totalResults, tt.mockReturn.err).
					Once()
			}

			req, err := testhelper.NewJSONRequest(http.MethodPut, requestUrl, tt.requestBody)
			if err != nil {
				t.Error("could not build request", err)
				return
			}

			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)

			testhelper.VerifyResponseWithMetadata(t, tt.want.responseCode, tt.want.responseParsed, resp)
		})
	}
}

func TestCreateNotification(t *testing.T) {
	someNotification := getNotification()
	someNotification.Id = "to be ignored"

	wantNotification := getNotification()
	wantNotification.Id = "new id"

	type mockServiceReturn struct {
		item models.Notification
		err  error
	}
	type want struct {
		notificationServiceArg *models.Notification
		responseCode           int
		responseParsed         query.ResponseWithMetadata[models.Notification]
	}

	tests := []struct {
		name                 string
		notificationToCreate models.Notification
		mockServiceReturn    mockServiceReturn
		want                 want
	}{
		{
			name:                 "services are called with the correct parameters (read only fields don't affect outcome)",
			notificationToCreate: someNotification,
			mockServiceReturn:    mockServiceReturn{item: wantNotification},
			want: want{
				notificationServiceArg: helper.ToPtr(someNotification),
				responseCode:           http.StatusCreated,
				responseParsed:         query.ResponseWithMetadata[models.Notification]{Data: wantNotification},
			},
		},
		{
			name:                 "return internal server error on service failure",
			notificationToCreate: someNotification,
			mockServiceReturn:    mockServiceReturn{item: models.Notification{}, err: errors.New("some internal error")},
			want: want{
				notificationServiceArg: helper.ToPtr(someNotification),
				responseCode:           http.StatusInternalServerError,
			},
		},
		{
			name:                 "don't create a notification if validation fails",
			notificationToCreate: models.Notification{}, // invalid: mandatory parameters not set
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}

	httpMethod := http.MethodPost
	requestUrl := "/notifications"

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotificationService := mocks.NewNotificationService(t)

			// Create a new engine for testing
			engine := gin.Default()
			// constructor registers the routes
			_ = NewNotificationController(&engine.RouterGroup, mockNotificationService, testhelper.MockAuthMiddleware)

			if tt.want.notificationServiceArg != nil {
				mockNotificationService.EXPECT().CreateNotification(mock.Anything, *tt.want.notificationServiceArg).
					Return(tt.mockServiceReturn.item, tt.mockServiceReturn.err).
					Once()
			}

			req, err := testhelper.NewJSONRequest(httpMethod, requestUrl, tt.notificationToCreate)
			if err != nil {
				t.Error("could not build request", err)
				return
			}

			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)

			testhelper.VerifyResponseWithMetadata(t, tt.want.responseCode, tt.want.responseParsed, resp)
		})
	}
}
