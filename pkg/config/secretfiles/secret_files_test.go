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
	secretFiles := map[string]string{
		"db_password":             "  db_password   \n\n\t",
		"encryption_key_password": "  enc_key_pw\n",
		"encryption_key_salt":     `  enc_key_salt\&*\$# `,
	}
	tempDir := t.TempDir()
	for file, content := range secretFiles {
		err := os.WriteFile(tempDir+"/"+file, []byte(content), 0644)
		require.NoError(t, err)
	}

	tests := map[string]struct {
		envVars     map[string]string
		inputConfig config.Config
		wantConfig  config.Config
		wantErr     bool
	}{
		"read all secrets from files": {
			inputConfig: config.Config{},
			envVars: map[string]string{
				"DB_PASSWORD_FILE":                           tempDir + "/db_password",
				"DATABASE_ENCRYPTION_KEY_PASSWORD_FILE":      tempDir + "/encryption_key_password",
				"DATABASE_ENCRYPTION_KEY_PASSWORD_SALT_FILE": tempDir + "/encryption_key_salt",
			},
			wantConfig: config.Config{
				Database: config.Database{
					Password: `db_password`,
				},
				DatabaseEncryptionKey: config.DatabaseEncryptionKey{
					Password:     `enc_key_pw`,
					PasswordSalt: `enc_key_salt\&*\$#`,
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
