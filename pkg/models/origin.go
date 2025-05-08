// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

// Origin of an event, uniquely identified by the combination of namespace and class.
type Origin struct {
	Name      string `json:"name" binding:"required"`
	Class     string `json:"class" binding:"required"`
	Namespace string `json:"namespace" readonly:"true"`
}

// OriginReference is the reference to an origin, uniquely identified by the combination of namespace and class.
// The name is simply informational.
type OriginReference struct {
	Name      string `json:"name" readonly:"true"`
	Class     string `json:"class" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}
