// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	_ "github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
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
	gc.Header(web.APIVersionKey, web.APIVersion)

	notification, err := parseAndValidateNotification(gc)
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not get notification", err)
		return
	}

	notificationNew, err := c.notificationService.CreateNotification(gc, notification)
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not create notification", err)
		return
	}

	gc.JSON(http.StatusCreated, query.ResponseWithMetadata[models.Notification]{Data: notificationNew})
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
	gc.Header(web.APIVersionKey, web.APIVersion)

	resultSelector, err := helper.PrepareResultSelector(gc, []filter.RequestOption{}, []string{}, helper.ResultSelectorDefaults())
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not prepare result selector", err)
		return
	}

	notifications, totalResults, err := c.notificationService.ListNotifications(gc, resultSelector)
	if err != nil {
		restErrorHandler.ErrorHandler(gc, "could not list notifications", err)
		return
	}

	gc.JSON(http.StatusOK, query.ResponseListWithMetadata[models.Notification]{
		Metadata: query.NewMetadata(resultSelector, totalResults),
		Data:     notifications,
	})
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
	gc.Header(web.APIVersionKey, web.APIVersion)

	// for now we don't support filtering
	permittedFilters := []query.FilterOption{}

	gc.JSON(http.StatusOK, query.ResponseWithMetadata[[]query.FilterOption]{Data: permittedFilters})
}

func parseAndValidateNotification(gc *gin.Context) (notification models.Notification, err error) { // TODO: refine
	var empty models.Notification
	err = gc.ShouldBindJSON(&notification)
	if err != nil {
		return empty, &errs.ErrValidation{Message: fmt.Sprintf("can't parse body: %v", err)}
	}

	return notification, nil
}
