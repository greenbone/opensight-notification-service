// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"github.com/gin-gonic/gin"
	logsMiddleware "github.com/greenbone/opensight-golang-libraries/pkg/logs/ginMiddleware"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

func NewWebEngine(httpConfig config.Http, registry *errmap.Registry) *gin.Engine {
	ginWebEngine := gin.New()
	ginWebEngine.Use(
		logsMiddleware.Logging(),
		gin.Recovery(),
		middleware.CORS(httpConfig.AllowedOrigins),
		middleware.ErrorHandler(gin.ErrorTypeAny),
	)
	ginWebEngine.Use(middleware.InterpretErrors(gin.ErrorTypePrivate, registry))
	return ginWebEngine
}
