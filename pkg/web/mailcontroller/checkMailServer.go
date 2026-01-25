package mailcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

	group.POST("/check", ctrl.CheckMailServer)

	return ctrl
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

	err := mc.notificationChannelServicer.CheckNotificationChannelConnectivity(context.Background(), mailServer.ToModel())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
