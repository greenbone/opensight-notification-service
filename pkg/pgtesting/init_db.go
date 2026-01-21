// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package pgtesting

import (
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register the "postgres" driver
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
)

const driverName = "postgres" // depends on registered driver via import above

// NewDB is a helper that returns an open connection to a unique and isolated
// test database, fully migrated and ready to query.
func NewDB(t *testing.T) *sqlx.DB {
	t.Parallel() // each test has its own isolated database
	t.Helper()
	conf := pgtestdb.Config{ // must match the deployment in `compose.yml`
		DriverName: driverName,
		User:       "postgres",
		Password:   "password",
		Host:       "localhost",
		Port:       "5432",
		Options:    "sslmode=disable",
	}

	migrator := golangmigrator.New(
		repository.MigrationDir,
		golangmigrator.WithFS(repository.MigrationsFS),
	)
	db := pgtestdb.New(t, conf, migrator)
	return sqlx.NewDb(db, driverName) // wrap in sqlx as it is expected by the objects to test
}
