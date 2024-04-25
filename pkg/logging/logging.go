// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package logging

import (
	"fmt"

	"github.com/rs/zerolog"
)

const AllowedLogLevels string = "disabled|trace|debug|info|warn|error|fatal|panic"

func SetupLogger(logLevel string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed parsing log level: %w; only allowed levels are: %s", err, AllowedLogLevels)
	}

	zerolog.SetGlobalLevel(level)
	return nil
}
