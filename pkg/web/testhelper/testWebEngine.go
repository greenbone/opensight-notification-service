// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package testhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

func NewTestWebEngine(registry *errmap.Registry) *gin.Engine {
	gin.SetMode(gin.TestMode)

	ginWebEngine := gin.New()
	ginWebEngine.Use(
		gin.Recovery(),
		middleware.CORS([]string{
			"http://example.com",
		}),
		middleware.ErrorHandler(gin.ErrorTypeAny),
		middleware.InterpretErrors(gin.ErrorTypePrivate, registry),
	)

	return ginWebEngine
}
