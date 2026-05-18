// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getNotification() models.Notification {
	return models.Notification{
		Id:          "57fe22b8-89a4-445f-b6c7-ef9ea724ea48",
		Timestamp:   time.Time{}.Format(time.RFC3339Nano),
		Origin:      "Example Task XY",
		OriginClass: "serviceab/exampletaskxy",
		Title:       "Example Task XY failed",
		Detail:      "Example Task XY failed because ...",
		Level:       "error",
	}
}

func setup(t *testing.T) (*gin.Engine, *mocks.NotificationService) {
	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)
	mockNotificationService := mocks.NewNotificationService(t)
	AddNotificationController(router, mockNotificationService, authMiddleware)

	return router, mockNotificationService
}

func TestCreateNotification_ForbiddenRoles(t *testing.T) {
	t.Parallel()

	forbiddenRoles := []string{iam.OsiViewer, iam.User, iam.OsiUser, iam.OsiAdmin, iam.Admin, iam.NotificationAdmin}

	for _, role := range forbiddenRoles {
		t.Run("Create notification is forbidden for role "+role, func(t *testing.T) {
			t.Parallel()

			router, _ := setup(t)

			httpassert.New(t, router).
				Post("/notifications").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Expect().
				StatusCode(http.StatusForbidden)
		})
	}
}

func TestCreateNotification_AllowedRoles(t *testing.T) {
	t.Parallel()

	allowedRoles := []string{iam.Notification}

	for _, role := range allowedRoles {
		t.Run("Create notification is allowed for role "+role, func(t *testing.T) {
			t.Parallel()

			router, mockNotificationService := setup(t)

			mockNotificationService.EXPECT().CreateNotification(mock.Anything, mock.Anything).
				Maybe().
				Return(models.Notification{}, nil)

			httpassert.New(t, router).
				Post("/notifications").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				JsonContentObject(getNotification()).
				Expect().
				StatusCode(http.StatusCreated)
		})
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
			name:        "return bad request on invalid page size",
			requestBody: query.ResultSelector{Paging: &paging.Request{PageSize: -1}},
			want: want{
				serviceCall:  false,
				responseCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockNotificationService := setup(t)

			if tt.want.serviceCall {
				mockNotificationService.EXPECT().ListNotifications(mock.Anything, tt.want.serviceArg).
					Return(tt.mockReturn.items, tt.mockReturn.totalResults, tt.mockReturn.err).
					Once()
			}

			var gotResponse query.ResponseListWithMetadata[models.Notification]
			httpassert.New(t, router).Put("/notifications").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.User)).
				JsonContentObject(tt.requestBody).
				Expect().
				StatusCode(tt.want.responseCode).
				GetJsonBodyObject(&gotResponse)
			require.Equal(t, tt.want.responseParsed, gotResponse)
		})
	}
}

func TestListNotifications_ForbiddenRoles(t *testing.T) {
	t.Parallel()

	forbiddenRoles := []string{iam.Admin, iam.NotificationAdmin}

	for _, role := range forbiddenRoles {
		t.Run("List notifications is forbidden for role "+role, func(t *testing.T) {
			t.Parallel()

			router, _ := setup(t)

			httpassert.New(t, router).
				Put("/notifications").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Content("{}").
				Expect().
				StatusCode(http.StatusForbidden)
		})
	}
}

func TestListNotifications_AllowedRoles(t *testing.T) {
	t.Parallel()

	allowedRoles := []string{iam.OsiViewer, iam.User, iam.OsiUser, iam.OsiAdmin}

	for _, role := range allowedRoles {
		t.Run("List notifications is allowed for role "+role, func(t *testing.T) {
			t.Parallel()

			router, mockNotificationService := setup(t)

			mockNotificationService.EXPECT().ListNotifications(mock.Anything, mock.Anything).
				Return([]models.Notification{}, uint64(0), nil).
				Once()

			httpassert.New(t, router).
				Put("/notifications").
				AuthJwt(integrationTests.CreateJwtTokenWithRole(role)).
				Content("{}").
				Expect().
				StatusCode(http.StatusOK)
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
				notificationServiceArg: new(someNotification),
				responseCode:           http.StatusCreated,
				responseParsed:         query.ResponseWithMetadata[models.Notification]{Data: wantNotification},
			},
		},
		{
			name:                 "return internal server error on service failure",
			notificationToCreate: someNotification,
			mockServiceReturn:    mockServiceReturn{item: models.Notification{}, err: errors.New("some internal error")},
			want: want{
				notificationServiceArg: new(someNotification),
				responseCode:           http.StatusInternalServerError,
			},
		},
		{
			name:                 "don't create a notification if mandatory parameters not set",
			notificationToCreate: models.Notification{},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}

	requestUrl := "/notifications"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockNotificationService := setup(t)

			if tt.want.notificationServiceArg != nil {
				mockNotificationService.EXPECT().CreateNotification(mock.Anything, *tt.want.notificationServiceArg).
					Return(tt.mockServiceReturn.item, tt.mockServiceReturn.err).
					Once()
			}

			var gotResponse query.ResponseWithMetadata[models.Notification]
			httpassert.New(t, router).Post(requestUrl).
				JsonContentObject(tt.notificationToCreate).
				AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Notification)).
				Expect().
				StatusCode(tt.want.responseCode).
				GetJsonBodyObject(&gotResponse)

			require.Equal(t, tt.want.responseParsed, gotResponse)
		})
	}
}
