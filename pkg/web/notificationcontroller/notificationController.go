// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationcontroller

import (
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
)

type NotificationController struct {
	notificationService notificationservice.NotificationService
}

func AddNotificationController(
	router gin.IRouter,
	notificationService notificationservice.NotificationService,
	auth gin.HandlerFunc,
) {
	ctrl := &NotificationController{
		notificationService: notificationService,
	}

	groupPath := "/notifications"

	router.Group(groupPath).Use(middleware.AuthorizeRoles(auth, iam.User)...).
		PUT("", ctrl.ListNotifications).
		GET("/options", ctrl.GetOptions)
	// only to be used by other backend services
	router.Group(groupPath).Use(middleware.AuthorizeRoles(auth, iam.Notification)...).
		POST("", ctrl.CreateNotification)
}

// CreateNotification
//
//	@Summary		Create Notification
//	@Description	Create a new notification. It will always be stored by the notification service and it will possibly also trigger actions like sending mails, depending on the cofigured rules.
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

	var notification models.Notification
	if !ginEx.BindAndValidateBody(gc, &notification) {
		return
	}

	notificationNew, err := c.notificationService.CreateNotification(gc, notification)
	if ginEx.AddError(gc, err) {
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

	resultSelector, err := helper.PrepareResultSelector(gc, NotificationsRequestOptions, AllowedNotificationsSortFields, helper.ResultSelectorDefaults(DefaultSortingRequest))
	if ginEx.AddError(gc, err) {
		return
	}

	notifications, totalResults, err := c.notificationService.ListNotifications(gc, resultSelector)
	if ginEx.AddError(gc, err) {
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
//	@Success		200	{object}	query.ResponseWithMetadata[[]query.FilterOption]
//	@Header			all	{string}	api-version	"API version"
//	@Router			/notifications/options [get]
func (c *NotificationController) GetOptions(gc *gin.Context) {
	gc.Header(web.APIVersionKey, web.APIVersion)

	requestOptions := lo.Map(NotificationsRequestOptions, web.ToFilterOption)
	response := query.ResponseWithMetadata[[]query.FilterOption]{
		Data: requestOptions,
	}
	gc.JSON(http.StatusOK, response)
}
