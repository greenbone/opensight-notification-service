// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package repository_test

import (
	"strconv"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/peterldowns/pgtestdb"
	"github.com/stretchr/testify/require"
)

func TestNewClientMigrationStateSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configuredSchema string
		wantSchema       string
	}{
		"search path default": {
			wantSchema: "application",
		},
		"explicit schema": {
			configuredSchema: "migration_state",
			wantSchema:       "migration_state",
		},
		"quoted schema": {
			configuredSchema: `migration "state".schema`,
			wantSchema:       `migration "state".schema`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Match pkg/pgtesting/compose.yml, but leave application migrations to NewClient.
			testCfg := pgtestdb.Custom(t, pgtestdb.Config{
				DriverName: "postgres",
				User:       "postgres",
				Password:   "password",
				Host:       "localhost",
				Port:       "9632",
				Options:    "sslmode=disable",
				// autoMigrate's driver retains a connection outside the pool.
				ForceTerminateConnections: true,
			}, pgtestdb.NoopMigrator{})

			db, err := sqlx.Connect("postgres", testCfg.URL())
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			if testCase.configuredSchema == "" {
				// fulfill precondition: existing schema in search path
				_, err = db.Exec("CREATE SCHEMA application")
				require.NoError(t, err)
				_, err = db.Exec("ALTER DATABASE " + pq.QuoteIdentifier(testCfg.Database) +
					" SET search_path TO application, public")
				require.NoError(t, err)
			} else {
				var schemaExists bool
				require.NoError(t, db.Get(&schemaExists,
					"SELECT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = $1)",
					testCase.configuredSchema))
				require.False(t, schemaExists, "schema should not exist before migration")
			}
			require.NoError(t, db.Close())

			port, err := strconv.Atoi(testCfg.Port)
			require.NoError(t, err)
			cfg := config.Database{
				Host:                 testCfg.Host,
				Port:                 port,
				User:                 testCfg.User,
				Password:             testCfg.Password,
				DBName:               testCfg.Database,
				SSLMode:              "disable",
				MigrationStateSchema: testCase.configuredSchema,
			}
			firstDB, err := repository.NewClient(cfg)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, firstDB.Close()) })
			var stateSchemas []string
			require.NoError(t, firstDB.Select(&stateSchemas,
				"SELECT table_schema FROM information_schema.tables WHERE table_name = 'schema_migrations'"))
			require.Equal(t, []string{testCase.wantSchema}, stateSchemas)
			var version uint
			var dirty bool
			require.NoError(t, firstDB.QueryRow("SELECT version, dirty FROM "+
				pq.QuoteIdentifier(testCase.wantSchema)+".schema_migrations").Scan(&version, &dirty))
			require.Greater(t, version, uint(0))
			require.False(t, dirty)
			require.NoError(t, firstDB.Close())

			// Check that reopening continues from the previous migration state.
			reopenedDB, err := repository.NewClient(cfg)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, reopenedDB.Close()) })
			require.NoError(t, reopenedDB.QueryRow("SELECT version, dirty FROM "+
				pq.QuoteIdentifier(testCase.wantSchema)+".schema_migrations").Scan(&version, &dirty))
			require.Greater(t, version, uint(0))
			require.False(t, dirty)
		})
	}
}
