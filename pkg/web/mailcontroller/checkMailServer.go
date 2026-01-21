package mailcontroller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/dtos"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type CheckMailServerController struct {
	notificationChannelServicer notificationchannelservice.NotificationChannelServicer
}

func AddCheckMailServerController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.NotificationChannelServicer,
	auth gin.HandlerFunc,
) *CheckMailServerController {
	ctrl := &CheckMailServerController{
		notificationChannelServicer: notificationChannelServicer,
	}

	group := router.Group("/notification-channel/mail").
		Use(middleware.AuthorizeRoles(auth, "admin")...)
	group.Use(validationErrorHandler(gin.ErrorTypePrivate))

	group.POST("/check", ctrl.CheckMailServer)

	return ctrl
}

// TODO: 21.01.2026 stolksdorf - move
func validationErrorHandler(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, errorValue := range c.Errors.ByType(errorType) {
			validateErrors := dtos.ValidateErrors{}
			if errors.As(errorValue, &validateErrors) {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponses.NewErrorValidationResponse("", "", validateErrors))
				return
			}
		}
	}
}

// CheckMailServer
//
//	@Summary		Check mail server
//	@Description	Check if a mail server is reachable
//	@Tags			mailserver
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailServerConfig	body		dtos.CheckMailServerRequest	true	"Mail server to check"
//	@Success		204 "Mail server reachable"
//	@Failure		400			{object}	map[string]string
//	@Failure		422 "Mail server error"
//	@Router			/notifications/mail/check [post]
func (mc *CheckMailServerController) CheckMailServer(c *gin.Context) {
	var mailServer dtos.CheckMailServerRequest
	if err := c.ShouldBindJSON(&mailServer); err != nil {
		_ = c.Error(err)
		return
	}
	if err := mailServer.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	err := mc.notificationChannelServicer.CheckNotificationChannelConnectivity(context.Background(), mailServer)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
