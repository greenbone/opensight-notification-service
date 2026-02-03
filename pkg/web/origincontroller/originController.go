// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package origincontroller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type OriginService interface {
	UpsertOrigins(ctx context.Context, namespace string, origins []entities.Origin) error
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
	crtl := &OriginController{
		originService: originService,
	}
	crtl.RegisterRoutes(router, auth)

	return crtl
}

func (c *OriginController) RegisterRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/origins").
		// TODO: remove user role here, only services are supposed to register origins
		Use(middleware.AuthorizeRoles(auth, middleware.UserRole, middleware.NotificationRole)...)
	group.PUT("/:namespace", c.RegisterOrigins)
}

// RegisterOrigins
//
//	@Summary		Register Origins
//	@Description	Registers a set of origins in the given namespace. Replaces the content of the namespace if it already existed. The origins can be ulitized to set trigger conditions for actions.
//	@Tags			origin
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			namespace	path		string			true	"namespace of the calling service, need to be unique among all services registering origins"
//	@Param			origins		body		[]models.Origin	true	"origins provided by the calling service"
//	@Success		200			{object}	query.ResponseWithMetadata[[]models.Origin]
//	@Header			all			{string}	api-version	"API version"
//	@Router			/origins/{namespace} [put]
func (c *OriginController) RegisterOrigins(gc *gin.Context) {
	gc.Header(web.APIVersionKey, web.APIVersion)

	namespace := gc.Param("namespace")
	var origins []models.Origin
	origins, err := parseAndValidateOrigins(gc)
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not parse origins", err)
		return
	}

	originsEntities := make([]entities.Origin, 0, len(origins))
	for _, origin := range origins {
		originsEntities = append(originsEntities, origin.ToEntity())
	}

	err = c.originService.UpsertOrigins(gc.Request.Context(), namespace, originsEntities)
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not upsert origins", err)
		return
	}

	gc.Status(http.StatusNoContent)
}

func parseAndValidateOrigins(gc *gin.Context) (origins []models.Origin, err error) {
	err = gc.ShouldBindJSON(&origins)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return nil, &errs.ErrValidation{Message: "body must not be empty"}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, &errs.ErrValidation{Message: "body is not valid json"}
		default:
			return nil, &errs.ErrValidation{Message: fmt.Sprintf("invalid input: %v", err)}
		}
	}

	return origins, nil
}
