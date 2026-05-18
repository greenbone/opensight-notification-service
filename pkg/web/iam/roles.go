// SPDX-FileCopyrightText: 2024-2025 Greenbone AG
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package iam

const (
	OsiViewer         = "osi.viewer"
	User              = "user" // Stays to ensure backwards compatibility
	OsiUser           = "osi.user"
	Admin             = "admin" // Stays to ensure backwards compatibility
	OsiAdmin          = "osi.admin"
	Notification      = "opensight_notification_role" // Only a service user
	NotificationAdmin = "notification.admin"          // Same as Admin
)
