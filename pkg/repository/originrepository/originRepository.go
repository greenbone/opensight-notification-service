// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/jmoiron/sqlx"
)

type OriginRepository struct {
	client *sqlx.DB
}

func NewOriginRepository(db *sqlx.DB) (*OriginRepository, error) {
	if db == nil {
		return nil, errors.New("nil db reference")
	}
	r := &OriginRepository{
		client: db,
	}
	return r, nil
}

// UpsertOrigins replaces all origins for the given namespace with the provided ones.
// Note: `origins.Namespace` is ignored, only the provided `namespace` parameter is used.
func (r *OriginRepository) UpsertOrigins(ctx context.Context, namespace string, origins []entities.Origin) (err error) {
	if namespace == "" {
		return errors.New("namespace must not be empty")
	}

	var originRows []originRow
	for _, o := range origins {
		originRows = append(originRows, toOriginRow(o, namespace))
	}

	tx, err := r.client.BeginTxx(ctx, nil) // replacement of existing entries must be atomic
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // note: rollback after successful commit is a no-op

	_, err = tx.Exec(deleteOriginsQuery, namespace)
	if err != nil {
		return fmt.Errorf("could not delete existing origins: %w", err)
	}

	if len(originRows) != 0 {
		_, err = tx.NamedExec(createOriginsQuery, originRows)
		if err != nil {
			return fmt.Errorf("could not insert origins: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil // so far we have no usecase to return the created origins
}

// ListOrigins returns all origins in the database.
// The origins are ordered by namespace and name.
func (r *OriginRepository) ListOrigins(ctx context.Context) ([]entities.Origin, error) {
	// Note: so far we don't have a usecase for filter or pagination
	// We usually need all of them and the total number is expected to be reasonably small.

	var originRows []originRow
	err := r.client.SelectContext(ctx, &originRows, listOriginsQuery)
	if err != nil {
		return nil, fmt.Errorf("could not get origins: %w", err)
	}

	origins := make([]entities.Origin, 0, len(originRows))
	for _, row := range originRows {
		origins = append(origins, row.toOriginEntity())
	}

	return origins, nil
}
