// SPDX-FileCopyrightText: 2025 Greenbone Networks GmbH <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/samber/lo"
)

const UserRole = "user"
const AdminRole = "admin"
const NotificationRole = "opensight_notification_role"

// AuthorizeRoles adds role-based authorization middleware to the provided handler.
// It first executes the provided authentication function, then checks if the user has one of the required roles.
func AuthorizeRoles(authFunc gin.HandlerFunc, roles ...string) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{authFunc}

	if len(roles) > 0 {
		handlers = append(handlers, func(c *gin.Context) {
			if len(c.Errors) > 0 {
				return
			}

			userContext, err := auth.GetUserContext(c)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if lo.None(roles, userContext.Roles) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Next()
		})
	}

	return handlers
}
