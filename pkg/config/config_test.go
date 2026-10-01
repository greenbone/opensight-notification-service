// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package config_test

import (
	"os"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
)

func TestDatabase_MigrationStateSchema(t *testing.T) {
	for _, tt := range []struct {
		name  string
		set   bool
		value string
	}{
		{name: "unset"},
		{name: "empty", set: true},
		{name: "custom", set: true, value: "Migration-State\"Schema"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Setenv restores the original environment even for the unset case.
			t.Setenv("DB_MIGRATION_STATE_SCHEMA", tt.value)
			if !tt.set {
				require.NoError(t, os.Unsetenv("DB_MIGRATION_STATE_SCHEMA"))
			}
			var cfg config.Config
			require.NoError(t, envconfig.Process("", &cfg))
			require.Equal(t, tt.value, cfg.Database.MigrationStateSchema)
		})
	}
}
