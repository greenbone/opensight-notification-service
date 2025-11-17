// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package dtos

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
	"github.com/greenbone/opensight-notification-service/pkg/web"
)

var NotificationsRequestOptions = []filter.RequestOption{
	{
		Name: web.NameFilterRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeString,
		},
		Operators: filter.SortedReadableValues(
			web.OperatorEqual,
			web.OperatorContains,
			web.OperatorBeginsWith,
		),
		MultiSelect: true,
	},
	{
		Name: web.DescriptionFilterRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeString,
		},
		Operators: filter.SortedReadableValues(
			web.OperatorEqual,
			web.OperatorContains,
			web.OperatorBeginsWith,
		),
		MultiSelect: true,
	},
	{
		Name: web.OriginFilterRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeString,
		},
		Operators: filter.SortedReadableValues(
			web.OperatorEqual,
			web.OperatorContains,
			web.OperatorBeginsWith,
		),
		MultiSelect: true,
	},
	{
		Name: web.OccurrenceFilterRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeDateTime,
		},
		Operators: filter.SortedReadableValues(
			web.OperatorBefore,
			web.OperatorAfter,
		),
	},
	{
		Name: web.LevelFilterRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeEnum,
		},
		Operators: filter.SortedReadableValues(
			web.OperatorEqual,
		),
		Values:      []string{"info", "warning", "error"},
		MultiSelect: true,
	},
}

var AllowedNotificationsSortFields = []string{dtos.NameField, dtos.DescriptionFieldName, dtos.OccurrenceFieldName, dtos.LevelFieldName, dtos.OriginFieldName}

var DefaultSortingRequest = &sorting.Request{
	SortColumn:    dtos.OccurrenceFieldName, // default sort by latest notification
	SortDirection: sorting.DirectionDescending,
}
