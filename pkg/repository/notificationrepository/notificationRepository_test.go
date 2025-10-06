// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"context"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateNotification_ListNotification(t *testing.T) {
	resultSelectorListAll := query.ResultSelector{
		Paging: &paging.Request{
			PageSize: 1000,
		},
	}

	tests := map[string]struct {
		notificationIn   models.Notification
		wantNotification models.Notification
		wantErr          bool
	}{
		"create notification": {
			notificationIn: models.Notification{
				Id:        "read only, to be ignored",
				Origin:    "test",
				OriginUri: "vi/test",
				Timestamp: "2024-10-10T10:00:00Z",
				Title:     "Test Notification",
				Detail:    "This is a test notification",
				Level:     "info",
				CustomFields: map[string]any{
					"key1": "value1",
					"key2": int(2),
					"key3": []string{"a", "b", "c"},
				},
			},
			wantNotification: models.Notification{
				Id:        "", // will be set after creation by db
				Origin:    "test",
				OriginUri: "vi/test",
				Timestamp: "2024-10-10T10:00:00Z",
				Title:     "Test Notification",
				Detail:    "This is a test notification",
				Level:     "info",
				CustomFields: map[string]any{
					"key1": "value1",
					"key2": float64(2), // not all types are preserved by json marshal/unmarshal
					"key3": []any{"a", "b", "c"},
				},
			},
		},
		"fail on inserting invalid notification": {
			notificationIn: models.Notification{}, // missing required fields
			wantErr:        true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := pgtesting.NewDB(t)

			repo, err := NewNotificationRepository(db)
			require.NoError(t, err)

			ctx := context.Background()

			gotNotification, err := repo.CreateNotification(ctx, tt.notificationIn)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, gotNotification.Id)
				tt.wantNotification.Id = gotNotification.Id // set the ID for comparison
				assert.Equal(t, tt.wantNotification, gotNotification)

				fetchedNotifications, gotTotalResults, err := repo.ListNotifications(ctx, resultSelectorListAll)
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), gotTotalResults, "did not get expected number of results")
				assert.Len(t, fetchedNotifications, 1)
				assert.Equal(t, tt.wantNotification, fetchedNotifications[0])
			}
		})
	}
}

func Test_ListNotifications(t *testing.T) {
	db := pgtesting.NewDB(t)

	repo, err := NewNotificationRepository(db)
	require.NoError(t, err)

	notification1 := models.Notification{
		Origin:    "test",
		OriginUri: "vi/test",
		Timestamp: "2024-10-10T10:00:00Z",
		Title:     "Test Notification",
		Detail:    "This is a test notification",
		Level:     "info",
		CustomFields: map[string]any{
			"key1": "value1",
			"key2": int(2),
			"key3": []string{"a", "b", "c"},
		},
	}
	notification2 := models.Notification{
		Origin:    "test",
		OriginUri: "vi/test",
		Timestamp: "2024-11-10T10:00:00Z",
		Title:     "Test Notification 2",
		Detail:    "This is a second test notification",
		Level:     "error",
	}
	notification3 := models.Notification{
		Origin:    "test2",
		OriginUri: "vi/test2",
		Timestamp: "2024-12-10T10:00:00Z",
		Title:     "Test Notification 3",
		Detail:    "This is a third test notification",
		Level:     "warning",
	}

	wantNotification1 := notification1
	wantNotification1.CustomFields = map[string]any{
		"key1": "value1",
		"key2": float64(2), // not all types are preserved by json marshal/unmarshal
		"key3": []any{"a", "b", "c"},
	}
	wantNotification2 := notification2
	wantNotification3 := notification3

	// sufficiently large page size to get all results in one page
	var bigPage = &paging.Request{
		PageSize: 1000,
	}

	notifications := []models.Notification{notification1, notification2, notification3}
	wantNotifications := []*models.Notification{&wantNotification1, &wantNotification2, &wantNotification3}

	ctx := context.Background()
	for ii, notification := range notifications {
		createdNotification, err := repo.CreateNotification(ctx, notification)
		require.NoError(t, err)
		require.NotEmpty(t, createdNotification.Id)
		wantNotifications[ii].Id = createdNotification.Id // set the ID for comparison
	}

	tests := map[string]struct {
		resultSelector    query.ResultSelector
		wantNotifications []models.Notification
		wantTotalResults  uint64
	}{
		"all notifications": {
			resultSelector: query.ResultSelector{
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{
				wantNotification1,
				wantNotification2,
				wantNotification3,
			},
			wantTotalResults: 3,
		},
		"results on several pages, first page": {
			resultSelector: query.ResultSelector{
				Paging: &paging.Request{
					PageSize: 2,
				},
				Sorting: &sorting.Request{ // tests also sorting
					SortColumn:    dtos.OccurrenceFieldName,
					SortDirection: sorting.DirectionAscending, // oldest first
				},
			},
			wantNotifications: []models.Notification{
				wantNotification1,
				wantNotification2,
			},
			wantTotalResults: 3, // total results independent of page size
		},
		"results on several pages, second page": {
			resultSelector: query.ResultSelector{
				Paging: &paging.Request{
					PageSize:  2,
					PageIndex: 1,
				},
				Sorting: &sorting.Request{
					SortColumn:    dtos.OccurrenceFieldName,
					SortDirection: sorting.DirectionAscending, // oldest first
				},
			},
			wantNotifications: []models.Notification{
				wantNotification3,
			},
			wantTotalResults: 3, // total results independent of page size
		},
		"filtery by name": {
			resultSelector: query.ResultSelector{
				Filter: &filter.Request{
					Operator: filter.LogicOperatorAnd,
					Fields: []filter.RequestField{
						{
							Name:     dtos.NameField,
							Value:    "Test Notification 2",
							Operator: filter.CompareOperatorIsEqualTo,
						},
					},
				},
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{
				wantNotification2,
			},
			wantTotalResults: 1,
		},
		"filter by origin": {
			resultSelector: query.ResultSelector{
				Filter: &filter.Request{
					Operator: filter.LogicOperatorAnd,
					Fields: []filter.RequestField{
						{
							Name:     dtos.OriginFieldName,
							Value:    "test",
							Operator: filter.CompareOperatorIsEqualTo,
						},
					},
				},
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{
				wantNotification1,
				wantNotification2,
			},
			wantTotalResults: 2,
		},
		"filter by level": {
			resultSelector: query.ResultSelector{
				Filter: &filter.Request{
					Operator: filter.LogicOperatorAnd,
					Fields: []filter.RequestField{
						{
							Name:     dtos.LevelFieldName,
							Value:    "warning",
							Operator: filter.CompareOperatorIsEqualTo,
						},
					},
				},
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{
				wantNotification3,
			},
			wantTotalResults: 1,
		},
		"filter by occurence": {
			resultSelector: query.ResultSelector{
				Filter: &filter.Request{
					Operator: filter.LogicOperatorAnd,
					Fields: []filter.RequestField{
						{
							Name:     dtos.OccurrenceFieldName,
							Value:    "2024-11-10T10:00:00Z",
							Operator: filter.CompareOperatorIsEqualTo,
						},
					},
				},
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{
				wantNotification2,
			},
			wantTotalResults: 1,
		},
		"no results": {
			resultSelector: query.ResultSelector{
				Filter: &filter.Request{
					Operator: filter.LogicOperatorAnd,
					Fields: []filter.RequestField{
						{
							Name:     dtos.OccurrenceFieldName,
							Value:    "1000-10-03T12:00:00Z",
							Operator: filter.CompareOperatorIsEqualTo,
						},
					},
				},
				Paging: bigPage,
			},
			wantNotifications: []models.Notification{},
			wantTotalResults:  0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotScans, totalResults, err := repo.ListNotifications(ctx, tt.resultSelector)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTotalResults, totalResults)
			assert.ElementsMatch(t, tt.wantNotifications, gotScans)
		})
	}
}
