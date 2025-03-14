// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package swagger

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/greenbone/opensight-golang-libraries/pkg/swagger"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/swaggo/swag"
)

func GetApiDocsHandler(specs *swag.Spec, config config.KeycloakConfig) gin.HandlerFunc {
	authConfig := &ginSwagger.OAuthConfig{
		ClientId: config.WebClientName,
		// other config values are ignored for now and are already set by the caller beforehand
	}

	apiDocsHandler := ginSwagger.GinWrapHandler(
		ginSwagger.InstanceName(specs.InstanceName()),
		ginSwagger.OAuth(authConfig),
	)
	return apiDocsHandler
}
