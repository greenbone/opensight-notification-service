// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package repository

import (
	"fmt"
	pgQuery "github.com/greenbone/opensight-golang-libraries/pkg/postgres/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
)

func BuildQuery(resultSelector query.ResultSelector, baseQuery string, fieldMapping map[string]string) (string, []any, error) {
	querySettings := &pgQuery.Settings{
		FilterFieldMapping: fieldMapping,
	}

	qb := pgQuery.NewPostgresQueryBuilder(querySettings)
	queryCondition, arg, err := qb.Build(resultSelector)
	if err != nil {
		return "", nil, fmt.Errorf("error building query condition: %w", err)
	}

	listQuery := baseQuery + ` ` + queryCondition

	return listQuery, arg, nil
}
