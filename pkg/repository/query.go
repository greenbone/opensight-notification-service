// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package repository

import (
	"fmt"
	pgQuery "github.com/greenbone/opensight-golang-libraries/pkg/postgres/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
)

// helper function to build the query condition
func buildQueryCondition(resultSelector query.ResultSelector, baseQuery string, fieldMapping map[string]string) (string, []any, error) {
	querySettings := &pgQuery.Settings{
		FilterFieldMapping: fieldMapping,
	}

	qb := pgQuery.NewPostgresQueryBuilder(querySettings)
	queryCondition, args, err := qb.Build(resultSelector)
	if err != nil {
		return "", nil, fmt.Errorf("error building query condition: %w", err)
	}

	fullQuery := baseQuery + ` ` + queryCondition
	return fullQuery, args, nil
}

// BuildQuery builds a query for retrieving results based on the provided result selector
func BuildQuery(resultSelector query.ResultSelector, baseQuery string, fieldMapping map[string]string) (string, []any, error) {
	return buildQueryCondition(resultSelector, baseQuery, fieldMapping)
}

// BuildCountQuery builds a count query based on the provided filter request
func BuildCountQuery(filterRequest *filter.Request, baseQuery string, fieldMapping map[string]string) (string, []any, error) {
	if filterRequest == nil {
		return baseQuery, nil, nil
	}

	resultSelector := query.ResultSelector{
		Filter: filterRequest,
	}

	return buildQueryCondition(resultSelector, baseQuery, fieldMapping)
}
