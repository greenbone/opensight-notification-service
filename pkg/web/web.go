// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

func NewWebEngine(httpConfig config.Http) *gin.Engine {
	ginWebEngine := gin.New()
	ginWebEngine.Use(
		middleware.Logging(),
		gin.Recovery(),
		middleware.CORS(httpConfig.AllowedOrigins),
		middleware.ErrorHandler(gin.ErrorTypeAny),
	)
	return ginWebEngine
}
