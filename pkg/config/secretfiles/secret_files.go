// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package secretfiles

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/secretfiles"
	"github.com/greenbone/opensight-notification-service/pkg/config"
)

const (
	dbPasswordPathEnvVar = "DB_PASSWORD_FILE"
)

// Read takes the filepaths from environment variables and parses the content
// into the respective secret inside the passed config.
// A failure can have side effects on the passed config, so error from this function
// should be treated as fatal.
func Read(cfg *config.Config) (err error) {
	return secretfiles.ReadSecret(dbPasswordPathEnvVar, &cfg.Database.Password)
}
