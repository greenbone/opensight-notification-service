// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package dtos

type VersionResponseDto struct {
	Version string `json:"version" example:"0.0.1-alpha1-dev1"`
}
