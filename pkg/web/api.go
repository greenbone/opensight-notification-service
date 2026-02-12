// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package web

import (
	"strings"

	"github.com/gin-gonic/gin"
	docs "github.com/greenbone/opensight-notification-service/api/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/swagger"
)

// comment block for api docs generation via swag:

//	@securitydefinitions.oauth2.implicit	KeycloakAuth
//	@authorizationUrl						{{.KeycloakAuthUrl}}/realms/{{.KeycloakRealm}}/protocol/openid-connect/auth

//	@title			Notification Service API
//	@version		1.0
//	@description	HTTP API of the Notification service

//	@license.name	AGPL-3.0-or-later

//	@BasePath	/api/notification-service

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

const (
	// APIVersion is the current version of the web API
	APIVersion = "1.0"
)

const APIVersionKey = "api-version"

func RegisterSwaggerDocsRoute(docsRouter gin.IRouter, kc config.KeycloakConfig) {
	docs.SwaggerInfonotificationservice.SwaggerTemplate = strings.ReplaceAll(docs.SwaggerInfonotificationservice.SwaggerTemplate,
		"{{.KeycloakAuthUrl}}", kc.PublicUrl)
	docs.SwaggerInfonotificationservice.SwaggerTemplate = strings.ReplaceAll(docs.SwaggerInfonotificationservice.SwaggerTemplate,
		"{{.KeycloakRealm}}", kc.Realm)
	apiDocsHandler := swagger.GetApiDocsHandler(docs.SwaggerInfonotificationservice, kc)
	docsRouter.GET("/notification-service/*any", apiDocsHandler)
}
