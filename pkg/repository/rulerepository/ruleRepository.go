// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulerepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const pgErrorUniqueViolationCode = "23505"

var ErrInvalidID error = errors.New("id is not a valid uuid-v4")
var ErrDuplicateRuleName error = errors.New("rule with the same name already exists")

type RuleRepository struct {
	client *sqlx.DB
}

func NewRuleRepository(db *sqlx.DB) (*RuleRepository, error) {
	if db == nil {
		return nil, errors.New("nil db reference")
	}
	r := &RuleRepository{
		client: db,
	}
	return r, nil
}

func (r *RuleRepository) Get(ctx context.Context, id string) (models.Rule, error) {
	err := validateId(id)
	if err != nil {
		return models.Rule{}, err
	}

	var row ruleRow
	if err := r.client.GetContext(ctx, &row, getRuleByIdQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Rule{}, errs.ErrItemNotFound
		}
		return models.Rule{}, fmt.Errorf("select by id failed: %w", err)
	}

	return row.ToModel()
}

func (r *RuleRepository) List(ctx context.Context) ([]models.Rule, error) {
	var rows []ruleRow
	if err := r.client.SelectContext(ctx, &rows, listRulesUnfilteredQuery); err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	var rules []models.Rule
	for _, row := range rows {
		ruleModel, err := row.ToModel()
		if err != nil {
			return nil, fmt.Errorf("failed to convert rule to model: %w", err)
		}
		rules = append(rules, ruleModel)
	}

	return rules, nil
}

func (r *RuleRepository) Create(ctx context.Context, rule models.Rule) (models.Rule, error) {
	rowIn := toRuleRow(rule)

	createStatement, err := r.client.PrepareNamedContext(ctx, createRuleQuery)
	if err != nil {
		return models.Rule{}, fmt.Errorf("could not prepare sql statement: %w", err)
	}

	var row ruleRow
	err = createStatement.QueryRowxContext(ctx, rowIn).StructScan(&row)
	if err != nil {
		err = postgresErrorHandling(err)
		return models.Rule{}, fmt.Errorf("could not create rule: %w", err)
	}

	return row.ToModel()
}

func (r *RuleRepository) Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error) {
	err := validateId(id)
	if err != nil {
		return models.Rule{}, err
	}

	rowIn := toRuleRow(rule)

	updateStatement, err := r.client.PrepareNamedContext(ctx, updateRuleQuery)
	if err != nil {
		return models.Rule{}, fmt.Errorf("could not prepare sql statement: %w", err)
	}

	var row ruleRow
	rowIn.ID = id
	err = updateStatement.QueryRowxContext(ctx, rowIn).StructScan(&row)
	if err != nil {
		err = postgresErrorHandling(err)
		return models.Rule{}, fmt.Errorf("could not update rule: %w", err)
	}

	return row.ToModel()
}

func (r *RuleRepository) Delete(ctx context.Context, id string) error {
	err := validateId(id)
	if err != nil {
		return err
	}

	_, err = r.client.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}

// ValidateId verifies if the passed id is a valid uuid-v4
// Note: the id format is enforced by the table definitions; we should validate before executing the query to avoid unnecessary calls.
// Furthermore there is no easy way to extract this error information from the postgres driver error without revealing implementation details.
func validateId(id string) error {
	err := validation.Validate.Var(id, "uuid4")
	if err != nil {
		return ErrInvalidID
	}

	return nil
}

func postgresErrorHandling(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrItemNotFound
	}

	if pgErr, ok := errors.AsType[*pq.Error](err); ok {
		if pgErr.Code == pgErrorUniqueViolationCode {
			return ErrDuplicateRuleName
		}
	}

	return fmt.Errorf("error querying database: %w", err)
}
