// SPDX-FileCopyrightText: 2023 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
)

type NotificationController struct {
	notificationService port.NotificationService
}

func NewNotificationController(router gin.IRouter, notificationService port.NotificationService) *NotificationController {
	ctrl := &NotificationController{
		notificationService: notificationService,
	}

	ctrl.registerRoutes(router)

	return ctrl
}

func (c *NotificationController) registerRoutes(router gin.IRouter) {
	group := router.Group("/notifications")
	group.POST("", c.CreateNotification)
	group.PUT("", c.ListNotifications)
	group.GET("/options", c.GetOptions)
}

// CreateNotification
//
//	@Summary		Create Notification
//	@Description	Create a new notification
//	@Tags			notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Authentication header"	example(Bearer eyJhbGciOiJSUzI1NiIs)
//	@Param			Notification	body		models.Notification	true	"notification to add"
//	@Success		201				{object}	query.ResponseWithMetadata[models.Notification]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/notifications [post]
func (c *NotificationController) CreateNotification(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// ListNotifications
//
//	@Summary		List Notifications
//	@Description	Returns a list of notifications matching the provided filters
//	@Tags			notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Authentication header"	example(Bearer eyJhbGciOiJSUzI1NiIs)
//	@Param			MatchCriterias	body		query.ResultSelector	true	"filters, paging and sorting"
//	@Success		200				{object}	query.ResponseListWithMetadata[models.Notification]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/notifications [put]
func (c *NotificationController) ListNotifications(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}

// GetOptions
//
//	@Summary		Notification filter options
//	@Description	Get filter options for listing notifications
//	@Tags			notification
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authentication header"	example(Bearer eyJhbGciOiJSUzI1NiIs)
//	@Success		200				{object}	query.ResponseWithMetadata[[]query.FilterOption]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/notifications/options [get]
func (c *NotificationController) GetOptions(gc *gin.Context) {
	gc.Status(http.StatusNotImplemented)
}
