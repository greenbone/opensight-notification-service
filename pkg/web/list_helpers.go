// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/dtos"
)

func ToFilterOption(options filter.RequestOption, _ int) query.FilterOption {
	return query.FilterOption(options)
}

// Filter request options operator labels
const (
	LabelIsEqualTo  = "is equal to"
	LabelContains   = "contains"
	LabelBeginsWith = "begins with"
	LabelAfter      = "after"
	LabelBefore     = "before"
)

// Default labels for filter options
const (
	labelNameFilterField        = "Name"
	labelDescriptionFilterField = "Description"
	labelOccurrenceFilterField  = "Occurrence"
	labelLevelFilterField       = "Level"
	labelOriginFilterField      = "Origin"
)

// Compare operators
var (
	OperatorEqual = filter.ReadableValue[filter.CompareOperator]{
		Label: LabelIsEqualTo,
		Value: filter.CompareOperatorIsEqualTo,
	}
	OperatorContains = filter.ReadableValue[filter.CompareOperator]{
		Label: LabelContains,
		Value: filter.CompareOperatorContains,
	}
	OperatorBeginsWith = filter.ReadableValue[filter.CompareOperator]{
		Label: LabelBeginsWith,
		Value: filter.CompareOperatorBeginsWith,
	}
	OperatorAfter = filter.ReadableValue[filter.CompareOperator]{
		Label: LabelAfter,
		Value: filter.CompareOperatorAfterDate,
	}
	OperatorBefore = filter.ReadableValue[filter.CompareOperator]{
		Label: LabelBefore,
		Value: filter.CompareOperatorBeforeDate,
	}
)

func NewNameRequest(labelOverride string) filter.ReadableValue[string] {
	return newFilterRequestName(labelOverride, labelNameFilterField, dtos.NameField)
}

func NewDescriptionRequest(labelOverride string) filter.ReadableValue[string] {
	return newFilterRequestName(labelOverride, labelDescriptionFilterField, dtos.DescriptionFieldName)
}

func NewOccurrenceRequest(labelOverride string) filter.ReadableValue[string] {
	return newFilterRequestName(labelOverride, labelOccurrenceFilterField, dtos.OccurrenceFieldName)
}

func NewOriginRequest(labelOverride string) filter.ReadableValue[string] {
	return newFilterRequestName(labelOverride, labelOriginFilterField, dtos.OriginFieldName)
}

func NewLevelRequest(labelOverride string) filter.ReadableValue[string] {
	return newFilterRequestName(labelOverride, labelLevelFilterField, dtos.LevelFieldName)
}

func newFilterRequestName(labelOverride, labelDefault, value string) filter.ReadableValue[string] {
	label := labelDefault
	if labelOverride != "" {
		label = labelOverride
	}

	return filter.ReadableValue[string]{
		Label: label,
		Value: value,
	}
}
