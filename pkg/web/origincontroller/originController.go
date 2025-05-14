// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
)

type OriginController struct{}

// RegisterOrigins
//
//	@Summary		Register Origins
//	@Description	Registers a set of origins in the given namespace. Replaces the content of the namespace if it already existed. The origins can be ulitized to set trigger conditions for actions.
//	@Tags			origin
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth[eventprovider]
//	@Param			namespace	path		string			true	"namespace of the calling service, need to be unique among all services registering origins"
//	@Param			origins		body		[]models.Origin	true	"origins provided by the calling service"
//	@Success		200			{object}	query.ResponseWithMetadata[[]models.Origin]
//	@Header			all			{string}	api-version	"API version"
//	@Router			/origins/{namespace} [put]
func (c *OriginController) RegisterOrigins(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}
