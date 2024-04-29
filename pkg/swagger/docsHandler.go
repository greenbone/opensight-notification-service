// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

func GetApiDocsHandler(specs *swag.Spec) gin.HandlerFunc {
	apiDocsHandler := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		func(config *ginSwagger.Config) {
			config.InstanceName = specs.InstanceName() // use default config except for instance name
		},
	)
	return apiDocsHandler
}
