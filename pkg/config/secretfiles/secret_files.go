// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package secretfiles

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/secretfiles"
	"github.com/greenbone/opensight-notification-service/pkg/config"
)

const (
	dbPasswordPathEnvVar                  = "DB_PASSWORD_FILE"
	dbEncryptionKeyPasswordPathEnvVar     = "DATABASE_ENCRYPTION_KEY_PASSWORD_FILE"
	dbEncryptionKeyPasswordSaltPathEnvVar = "DATABASE_ENCRYPTION_KEY_PASSWORD_SALT_FILE"
)

// Read takes the filepaths from environment variables and parses the content
// into the respective secret inside the passed config.
// A failure can have side effects on the passed config, so error from this function
// should be treated as fatal.
func Read(cfg *config.Config) (err error) {
	if err := secretfiles.ReadSecret(dbPasswordPathEnvVar, &cfg.Database.Password); err != nil {
		return err
	}
	if err := secretfiles.ReadSecret(dbEncryptionKeyPasswordPathEnvVar, &cfg.DatabaseEncryptionKey.Password); err != nil {
		return err
	}
	if err := secretfiles.ReadSecret(dbEncryptionKeyPasswordSaltPathEnvVar, &cfg.DatabaseEncryptionKey.PasswordSalt); err != nil {
		return err
	}
	return nil
}
