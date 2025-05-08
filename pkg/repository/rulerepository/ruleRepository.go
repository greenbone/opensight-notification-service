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
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const pgForeignKeyViolationCode = "23503"
const pgErrorUniqueViolationCode = "23505"

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
	query := ruleQuerySelect + ` WHERE r.id = $1 ` + ruleQueryGroupBy

	var row ruleRow
	if err := r.client.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Rule{}, fmt.Errorf("element with ID %s not found %w", id, errs.ErrItemNotFound)
		}
		return models.Rule{}, fmt.Errorf("select by id failed: %w", err)
	}

	return row.ToModel()
}

func (r *RuleRepository) List(ctx context.Context) ([]models.Rule, error) {
	query := ruleQuerySelect + ruleQueryGroupBy + ` ORDER BY r.name`

	var rows []ruleRow
	if err := r.client.SelectContext(ctx, &rows, query); err != nil {
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

	// Validate that all origins exist (we don't have foreign key constraints here)
	err := r.checkOriginsExist(ctx, rowIn.TriggerOrigins)
	if err != nil {
		return models.Rule{}, err
	}

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

	// Retrieve full rule with joined data
	return r.Get(ctx, row.ID)
}

func (r *RuleRepository) Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error) {
	rowIn := toRuleRow(rule)

	err := r.checkOriginsExist(ctx, rowIn.TriggerOrigins)
	if err != nil {
		return models.Rule{}, err
	}

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

	return r.Get(ctx, id)
}

func (r *RuleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}

func (r *RuleRepository) checkOriginsExist(ctx context.Context, origins []string) error {
	if len(origins) == 0 {
		return nil
	}

	var count int
	err := r.client.GetContext(
		ctx,
		&count,
		`SELECT COUNT(DISTINCT class) FROM notification_service.origins WHERE class = ANY($1)`,
		pq.Array(origins),
	)
	if err != nil {
		return fmt.Errorf("failed to validate origins: %w", err)
	}
	if count != len(origins) {
		return &errs.ErrValidation{Message: "invalid rule", Errors: map[string]string{"trigger.origins": "one or more origins do not exist"}}
	}

	return nil
}

func postgresErrorHandling(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrItemNotFound
	}
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch errCode := pgErr.Code; errCode {
		case pgForeignKeyViolationCode:
			return &errs.ErrValidation{Message: "invalid rule", Errors: map[string]string{"trigger.action.channel": "channel does not exist"}}
		case pgErrorUniqueViolationCode: // unique_violation
			return &errs.ErrConflict{Message: "rule with the same name already exists", Errors: map[string]string{"name": "must be unique"}}
		}
	}

	return fmt.Errorf("error querying database: %w", err)
}
