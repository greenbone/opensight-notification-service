// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package repository

import (
	"fmt"

	pgquery "github.com/greenbone/opensight-golang-libraries/pkg/postgres/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
)

// BuildListQuery builds a query for retrieving results based on the provided result selector
func BuildListQuery(resultSelector query.ResultSelector, baseQuery string, querySettings pgquery.Settings) (string, []any, error) {
	qb, err := pgquery.NewPostgresQueryBuilder(querySettings)
	if err != nil {
		return "", nil, fmt.Errorf("error creating query builder: %w", err)
	}

	queryCondition, args, err := qb.Build(resultSelector)
	if err != nil {
		return "", nil, fmt.Errorf("error building query condition: %w", err)
	}

	fullQuery := baseQuery + ` ` + queryCondition
	return fullQuery, args, nil
}

// BuildCountQuery builds a count query based on the provided filter request
func BuildCountQuery(filterRequest *filter.Request, baseQuery string, querySettings pgquery.Settings) (string, []any, error) {
	// create a resultSelector for the filter, sorting and paging are intentionally omitted here.
	// sorting does not affect the count, and paging (limiting or offsetting rows) is not necessary for counting.
	resultSelector := query.ResultSelector{
		Filter: filterRequest, // only the filter is applied to narrow down the count based on conditions.
	}

	countQuery, args, err := BuildListQuery(resultSelector, baseQuery, querySettings)
	if err != nil {
		return "", nil, fmt.Errorf("error building count query: %w", err)
	}
	countQuery = "SELECT count(*) FROM (" + countQuery + ")"

	return countQuery, args, nil
}
