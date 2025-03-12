// SPDX-FileCopyrightText: 2025 Greenbone Networks GmbH <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package middleware

import (
	"github.com/gin-gonic/gin"
)

const UserRole = "user"

// AuthorizeRoles adds role-based authorization middleware to the provided handler.
// It first executes the provided authentication function, then checks if the user has one of the required roles.
func AuthorizeRoles(authFunc gin.HandlerFunc, roles ...string) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{authFunc}

	if len(roles) > 0 {
		handlers = append(handlers, func(c *gin.Context) {
			if len(c.Errors) > 0 {
				return
			}

			// TODO: handle authorize user role
			// uncomment role check as soon as we can authenticate roles with notification-service-client account on keycloak

			//userContext, err := auth.GetUserContext(c)
			//if err != nil {
			//	_ = c.AbortWithError(http.StatusInternalServerError, err)
			//	return
			//}

			//if lo.None(roles, userContext.Roles) {
			//	c.AbortWithStatus(http.StatusForbidden)
			//	return
			//}

			c.Next()
		})
	}

	return handlers
}
