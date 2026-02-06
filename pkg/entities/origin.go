// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package entities

type Origin struct {
	Name      string
	Class     string
	ServiceID string // read-only
}
