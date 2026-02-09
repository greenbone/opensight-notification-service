// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type OriginService interface {
	UpsertOrigins(ctx context.Context, serviceID string, origins []entities.Origin) error
	ListOrigins(ctx context.Context) ([]entities.Origin, error)
}

type OriginController struct {
	originService OriginService
}

func NewOriginController(
	router gin.IRouter,
	originService OriginService,
	auth gin.HandlerFunc,
) *OriginController {
	ctrl := &OriginController{
		originService: originService,
	}
	ctrl.RegisterRoutes(router, auth)

	return ctrl
}

func (c *OriginController) RegisterRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/origins").
		Use(middleware.AuthorizeRoles(auth, middleware.NotificationRole)...)
	group.PUT("/:serviceID", c.RegisterOrigins)
}

// RegisterOrigins
//
//	@Summary		Register Origins
//	@Description	Registers a set of origins in the given service. Replaces origins of this service if they already existed. The origins can be ulitized to set trigger conditions for actions.
//	@Tags			origin
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			serviceID	path	string			true	"serviceID of the calling service, needs to be unique among all services registering origins"
//	@Param			origins		body	[]models.Origin	true	"origins provided by the calling service"
//	@Success		204			"origins registered"
//	@Failure		400			{object}	errorResponses.ErrorResponse
//	@Failure		500			{object}	errorResponses.ErrorResponse
//	@Header			all			{string}	api-version	"API version"
//	@Router			/origins/{serviceID} [put]
func (c *OriginController) RegisterOrigins(gc *gin.Context) {
	gc.Header(web.APIVersionKey, web.APIVersion)

	serviceID := gc.Param("serviceID")
	var origins models.OriginList
	if !ginEx.BindAndValidateBody(gc, &origins) {
		return
	}

	originsEntities := make([]entities.Origin, 0, len(origins))
	for _, origin := range origins {
		originsEntities = append(originsEntities, origin.ToEntity())
	}

	err := c.originService.UpsertOrigins(gc.Request.Context(), serviceID, originsEntities)
	if err != nil {
		ginEx.AddError(gc, err)
		return
	}

	gc.Status(http.StatusNoContent)
}
