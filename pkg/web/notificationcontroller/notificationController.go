// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"errors"
	"fmt"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/notificationcontroller/dtos"
	"github.com/samber/lo"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/greenbone/opensight-notification-service/pkg/restErrorHandler"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
)

type NotificationController struct {
	notificationService port.NotificationService
}

func NewNotificationController(router gin.IRouter, notificationService port.NotificationService, auth gin.HandlerFunc) *NotificationController {
	ctrl := &NotificationController{
		notificationService: notificationService,
	}

	ctrl.registerRoutes(router, auth)

	return ctrl
}

func (c *NotificationController) registerRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/notifications").Use(middleware.AuthorizeRoles(auth, middleware.UserRole, middleware.NotificationRole)...)
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
//	@Security		KeycloakAuth
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
//	@Security		KeycloakAuth
//	@Param			MatchCriterias	body		query.ResultSelector	true	"filters, paging and sorting"
//	@Success		200				{object}	query.ResponseListWithMetadata[models.Notification]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/notifications [put]
func (c *NotificationController) ListNotifications(gc *gin.Context) {
	gc.Header(web.APIVersionKey, web.APIVersion)

	resultSelector, err := helper.PrepareResultSelector(gc, dtos.NotificationsRequestOptions, dtos.AllowedNotificationsSortFields, helper.ResultSelectorDefaults(dtos.DefaultSortingRequest))
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
//	@Security		KeycloakAuth
//	@Success		200				{object}	query.ResponseWithMetadata[[]query.FilterOption]
//	@Header			all				{string}	api-version	"API version"
//	@Router			/notifications/options [get]
func (c *NotificationController) GetOptions(gc *gin.Context) {
	gc.Header(web.APIVersionKey, web.APIVersion)

	requestOptions := lo.Map(dtos.NotificationsRequestOptions, web.ToFilterOption)
	response := query.ResponseWithMetadata[[]query.FilterOption]{
		Data: requestOptions,
	}
	gc.JSON(http.StatusOK, response)
}

func parseAndValidateNotification(gc *gin.Context) (notification models.Notification, err error) {
	var empty models.Notification
	err = gc.ShouldBindJSON(&notification)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return empty, &errs.ErrValidation{Message: "body must not be empty"}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return empty, &errs.ErrValidation{Message: "body is not valid json"}
		default:
			return empty, &errs.ErrValidation{Message: fmt.Sprintf("invalid input: %v", err)}
		}
	}

	// validating the timestamp format via gin is rather clumsy, as it does not allow usage of the time layout constants and returns only a cryptic error message,
	// so instead we do it manually here
	_, err = time.Parse(time.RFC3339Nano, notification.Timestamp)
	if err != nil {
		return empty, &errs.ErrValidation{Message: fmt.Sprintf("invalid timestamp format: %v", err)}
	}

	return notification, nil
}
