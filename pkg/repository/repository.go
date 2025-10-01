// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package repository

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

//go:embed migrations
var MigrationsFS embed.FS

// directory within [MigrationsFS] where migration files are located
var MigrationDir = "migrations"

func NewClient(postgres config.Database) (*sqlx.DB, error) {
	// note: even though some parameters are part of the url path, [url.PathEscape] does not fit as it does not escape `:`
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=10",
		url.QueryEscape(postgres.User), url.QueryEscape(postgres.Password),
		url.QueryEscape(postgres.Host), postgres.Port, url.QueryEscape(postgres.DBName),
		url.QueryEscape(postgres.SSLMode))
	//connect to the db
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres database '%s:%d': %w", postgres.Host, postgres.Port, err)
	}

	if automigrateErr := autoMigrate(connectionString); automigrateErr != nil {
		if errors.Is(automigrateErr, migrate.ErrNoChange) { // TODO: handle errNoChange when migration file is unchanged on restart of the app on the same environment
			log.Debug().Msg("nothing to migrate")
			return db, nil
		}
		return nil, fmt.Errorf("error automigrating db: %w", automigrateErr)
	}

	return db, nil
}

func autoMigrate(connectionString string) error {
	log.Debug().Msg("starting database migration")
	//migrating postgres database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	databaseDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	sourceDriver, err := iofs.New(MigrationsFS, MigrationDir)
	if err != nil {
		return fmt.Errorf("could not read migration files: %w", err)
	}

	migration, err := migrate.NewWithInstance("embedded file system", sourceDriver, "postgres", databaseDriver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if migrateErr := migration.Up(); migrateErr != nil {
		return fmt.Errorf("migration error: %w", migrateErr)
	}
	log.Debug().Msg("database migration done")
	return nil
}
