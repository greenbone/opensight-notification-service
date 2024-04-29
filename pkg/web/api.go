// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"github.com/gin-gonic/gin"
	docs "github.com/greenbone/opensight-notification-service/api/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/swagger"
)

// comment block for api docs generation via swag:

//	@title			Notification Service API
//	@version		1.0
//	@description	HTTP API of the Notification service

//	@license.name	AGPL-3.0-or-later

//	@BasePath	/api/notification-service

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func RegisterSwaggerDocsRoute(docsRouter gin.IRouter) {
	apiDocsHandler := swagger.GetApiDocsHandler(docs.SwaggerInfonotificationservice)
	docsRouter.GET("/notification-service/*any", apiDocsHandler)
}
