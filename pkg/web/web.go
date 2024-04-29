// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func NewWebEngine() *gin.Engine {
	ginWebEngine := gin.New()
	ginWebEngine.Use(logger.SetLogger(), gin.Recovery())
	return ginWebEngine
}
