// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package healthcontroller

import (
	"github.com/greenbone/opensight-notification-service/pkg/web/healthcontroller/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/port"

	app "github.com/greenbone/opensight-notification-service"
)

type HealthController struct {
	healthService port.HealthService
}

func NewHealthController(router gin.IRouter, healthService port.HealthService) *HealthController {
	ctrl := HealthController{
		healthService: healthService,
	}
	ctrl.registerRoutes(router)
	return &ctrl
}

func (c *HealthController) registerRoutes(router gin.IRouter) {
	group := router.Group("/health")
	group.GET("started", c.Started)
	group.GET("alive", c.Alive)
	group.GET("ready", c.Ready)

	router.GET("/api/notification-service/version", c.readVersion)
}

//	@Summary		Service health status Started
//	@Description	Endpoint for 'started' health probes
//	@Tags			health
//	@Success		200	"Started"
//	@Router			/health/started [get]
func (c *HealthController) Started(gc *gin.Context) {
	gc.Status(http.StatusOK)
}

//	@Summary		Service health status Alive
//	@Description	Endpoint for 'alive' health probes
//	@Tags			health
//	@Success		200	"Alive"
//	@Router			/health/alive [get]
func (c *HealthController) Alive(gc *gin.Context) {
	gc.Status(http.StatusOK)
}

//	@Summary		Service health status Ready
//	@Description	Indicates if the service is ready to serve traffic
//	@Tags			health
//	@Success		200	"Ready"
//	@Failure		404	"Not ready"
//	@Router			/health/ready [get]
func (c *HealthController) Ready(gc *gin.Context) {
	if c.healthService.Ready(gc.Request.Context()) {
		gc.Status(http.StatusOK)
	} else {
		gc.Status(http.StatusNotFound)
	}
}

// readVersion
//
//	@Summary	Read API version
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	dtos.VersionResponseDto
//	@Router		/api/notification-service/version [get]
func (*HealthController) readVersion(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.VersionResponseDto{
		Version: app.Version,
	})
}
