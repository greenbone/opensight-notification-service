// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package healthcontroller

import (
	"github.com/gin-gonic/gin"
	docs "github.com/greenbone/opensight-notification-service/api/health"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/swagger"
)

// comment block for api docs generation via swag:

//	@title			Health API
//	@version		1.0
//	@description	HTTP API for live probes

//	@license.name	AGPL-3.0-or-later

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func RegisterSwaggerDocsRoute(docsRouter gin.IRouter, kc config.KeycloakConfig) {
	apiDocsHandler := swagger.GetApiDocsHandler(docs.SwaggerInfohealth, kc)
	docsRouter.GET("/health/*any", apiDocsHandler)
}
