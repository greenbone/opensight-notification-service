// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// see https://github.com/lib/pq/blob/3d613208bca2e74f2a20e04126ed30bcb5c4cc27/error.go#L78
const pgErrCodeConflict = "23505"

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

// UpsertOrigins replaces all origins for the given serviceID with the provided ones.
// Note: `origins.ServiceID` is ignored, only the provided `serviceID` parameter is used.
func (r *OriginRepository) UpsertOrigins(ctx context.Context, serviceID string, origins []entities.Origin) (err error) {
	if serviceID == "" {
		return errors.New("serviceID must not be empty")
	}

	var originRows []originRow
	for _, o := range origins {
		originRows = append(originRows, toOriginRow(o, serviceID))
	}

	tx, err := r.client.BeginTxx(ctx, nil) // replacement of existing entries must be atomic
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // note: rollback after successful commit is a no-op

	// acquire exclusive lock, as concurrent upserts for the same serviceID can result
	// extra data or unique constraint violation
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtext($1))", serviceID)
	if err != nil {
		return fmt.Errorf("could not acquire lock: %w", err)
	}

	_, err = tx.Exec(deleteOriginsQuery, serviceID)
	if err != nil {
		return fmt.Errorf("could not delete existing origins: %w", err)
	}

	if len(originRows) != 0 {
		_, err = tx.NamedExec(createOriginsQuery, originRows)
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) { // postgres specific error handling
				if pgErr.Code == pgErrCodeConflict {
					err = &errs.ErrConflict{Message: "duplicate origin class"}
				}
			}
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
// The origins are ordered by serviceID and name.
func (r *OriginRepository) ListOrigins(ctx context.Context) ([]entities.Origin, error) {
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
