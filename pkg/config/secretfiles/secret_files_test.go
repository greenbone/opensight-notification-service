// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package secretfiles

import (
	"os"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	// create files containing secrets
	tempDir := t.TempDir()
	err := os.WriteFile(tempDir+"/db_password", []byte("  db_password   \n\n\t"), 0644)
	require.NoError(t, err)

	tests := map[string]struct {
		envVars     map[string]string
		inputConfig config.Config
		wantConfig  config.Config
		wantErr     bool
	}{
		"read all secrets from files": {
			inputConfig: config.Config{},
			envVars: map[string]string{
				"DB_PASSWORD_FILE": tempDir + "/db_password",
			},
			wantConfig: config.Config{
				Database: config.Database{
					Password: `db_password`,
				},
			},
			wantErr: false,
		},
		"failure with invalid path": {
			inputConfig: config.Config{},
			envVars: map[string]string{
				"DB_PASSWORD_FILE": "/invalid/path",
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// set the environment variables
			for key, value := range tt.envVars {
				err := os.Setenv(key, value)
				require.NoError(t, err)
			}

			err := Read(&tt.inputConfig)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantConfig, tt.inputConfig)
			}
		})
	}
}
