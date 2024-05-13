// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package helper

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/greenbone/opensight-notification-service/pkg/errs"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	DefaultLimit uint64 = 50
)

// PrepareResultSelector converts the common query parameters of a List<X> endpoint to a ResultSelector.
// Parameters are verified where possible and defaults are set.
func PrepareResultSelector(gc *gin.Context, filterOptions []filter.RequestOption, allowedSortFields []string, defaults query.ResultSelector) (resultSelector query.ResultSelector, err error) {
	resultSelector = query.ResultSelector{}
	if err = gc.ShouldBindJSON(&resultSelector); err != nil {
		return resultSelector, &errs.ErrValidation{Message: fmt.Sprintf("can't parse body: %v", err)}
	}

	//apply defaults
	resultSelector = applyDefaults(resultSelector, defaults)
	if resultSelector.Sorting != nil { // TODO: implement sorting, and remove this if-branch
		resultSelector.Sorting = nil
		log.Warn().Msg("sorting not yet supported, dropping sort property from request")
	}

	err = validate(resultSelector, filterOptions, allowedSortFields)
	if err != nil {
		return resultSelector, &errs.ErrValidation{Message: fmt.Sprintf("error validating result selector %v", err)}
	}

	return resultSelector, nil
}

// ResultSelectorDefaults holds default result selectors
func ResultSelectorDefaults() query.ResultSelector {
	return query.ResultSelector{
		Paging: &paging.Request{
			PageIndex: 0,
			PageSize:  int(DefaultLimit),
		},
		Sorting: &sorting.Request{ // TODO: implement sorting, then we need a different default per endpoint instead of this dummy
			SortColumn:    "",
			SortDirection: "",
		},
	}
}

func applyDefaults(resultSelector query.ResultSelector, defaults query.ResultSelector) query.ResultSelector {
	if resultSelector.Paging == nil {
		resultSelector.Paging = defaults.Paging
	}
	if resultSelector.Sorting == nil {
		resultSelector.Sorting = defaults.Sorting
	} else {
		if resultSelector.Sorting.SortColumn == "" {
			resultSelector.Sorting.SortColumn = defaults.Sorting.SortColumn
			resultSelector.Sorting.SortDirection = defaults.Sorting.SortDirection
		}

		if resultSelector.Sorting.SortDirection == "" {
			resultSelector.Sorting.SortDirection = defaults.Sorting.SortDirection
		}
	}
	return resultSelector
}

func validate(resultSelector query.ResultSelector, filterOptions []filter.RequestOption, allowedSortFields []string) error {
	err := filter.ValidateFilter(resultSelector.Filter, filterOptions)
	if err != nil {
		return err
	}

	err = validateSorting(resultSelector.Sorting, allowedSortFields)
	if err != nil {
		return err
	}

	err = validatePaging(resultSelector.Paging)
	if err != nil {
		return err
	}

	return nil
}

func validateSorting(sortingRequest *sorting.Request, allowedSortFields []string) error {
	if sortingRequest == nil {
		return nil
	}

	err := sorting.ValidateSortingRequest(sortingRequest)
	if err != nil {
		return err
	}

	if !lo.Contains(allowedSortFields, sortingRequest.SortColumn) {
		return sorting.NewSortingError(fmt.Sprintf("%s is no valid sort column, possible values: %s",
			sortingRequest.SortColumn, strings.Join(allowedSortFields, ", ")))
	}

	return nil
}

func validatePaging(pagingRequest *paging.Request) error {
	if pagingRequest == nil {
		return nil
	}

	if pagingRequest.PageIndex < 0 {
		return paging.NewPagingError("%d is no valid page index, it must be >= 0", pagingRequest.PageIndex)
	}

	if pagingRequest.PageSize <= 0 {
		return paging.NewPagingError("%d is no valid page size, it must be > 0", pagingRequest.PageSize)
	}

	return nil
}
