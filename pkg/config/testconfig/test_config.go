// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package testconfig

// Env variables to control test suites execution
const (
	// RunAllGoEnv when set triggers execution of all go tests
	RunAllGoEnv = "TEST_ALL_GO"
	// RunPostgresEnv when set triggers postgres tests execution
	RunPostgresEnv = "TEST_POSTGRES"
)
