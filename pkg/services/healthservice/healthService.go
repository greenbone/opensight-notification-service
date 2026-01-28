// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package healthservice

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type HealthService interface {
	Ready(ctx context.Context) (ready bool)
}

type healthService struct {
	pgClient *sqlx.DB
}

func NewHealthService(pgClient *sqlx.DB) HealthService {
	return &healthService{
		pgClient: pgClient,
	}
}

// Ready indicates if the service is ready to serve traffic.
// Check that databases are up and ready to serve data
func (s *healthService) Ready(ctx context.Context) (ready bool) {
	// check postgres health
	err := s.pgClient.Ping()
	if err != nil {
		log.Debug().Msgf("error pinging postgres database %v", err)
		return false
	}

	return true
}
