package mailcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/mapper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
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
	ctrl.registerRoutes(router, auth)
	return ctrl
}

func (mc *CheckMailServerController) registerRoutes(router gin.IRouter, auth gin.HandlerFunc) {
	group := router.Group("/notification-channel/mail/check").
		Use(middleware.AuthorizeRoles(auth, "admin")...)

	group.POST("", mc.CheckMailServer)
}

// CheckMailServer
//
//	@Summary		Check mail server
//	@Description	Check if a mail server is reachable
//	@Tags			mailserver
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailServerConfig	body		models.MailNotificationChannel	true	"Mail server to check"
//	@Success		204 "Mail server reachable"
//	@Failure		400			{object}	map[string]string
//	@Failure		422 "Mail server error"
//	@Router			/notifications/mail [post]
func (mc *CheckMailServerController) CheckMailServer(c *gin.Context) {
	var channel models.MailNotificationChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notificationChannel := mapper.MapMailToNotificationChannel(channel)

	err := mc.notificationChannelServicer.CheckNotificationChannelConnectivity(context.Background(), notificationChannel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
