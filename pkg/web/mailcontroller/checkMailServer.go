package mailcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/ginEx"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller/maildto"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
)

type CheckMailServerController struct {
	notificationChannelServicer notificationchannelservice.MailChannelService
}

func AddCheckMailServerController(
	router gin.IRouter,
	notificationChannelServicer notificationchannelservice.MailChannelService,
	auth gin.HandlerFunc,
	registry errmap.ErrorRegistry,
) *CheckMailServerController {
	ctrl := &CheckMailServerController{
		notificationChannelServicer: notificationChannelServicer,
	}

	group := router.Group("/notification-channel/mail").
		Use(middleware.AuthorizeRoles(auth, "admin")...)

	group.POST("/check", ctrl.CheckMailServer)

	ctrl.configureMappings(registry)

	return ctrl
}

func (mc *CheckMailServerController) configureMappings(r errmap.ErrorRegistry) {
	r.Register(
		notificationchannelservice.ErrGetMailChannel,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
	r.Register(
		notificationchannelservice.ErrCreateMailClient,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Unable to create mail client"),
	)

	r.Register(
		notificationchannelservice.ErrMailServerUnreachable,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Server is unreachable"),
	)
}

// CheckMailServer
//
//	@Summary		Check mail server
//	@Description	Check if a mail server is reachable
//	@Tags			mailserver
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			MailServerConfig	body		maildto.CheckMailServerRequest	true	"Mail server to check"
//	@Success		204 "Mail server reachable"
//	@Failure		400			{object}	map[string]string
//	@Failure		422 "Mail server error"
//	@Router			/notifications/mail/check [post]
func (mc *CheckMailServerController) CheckMailServer(c *gin.Context) {
	var mailServer maildto.CheckMailServerRequest
	if !ginEx.BindAndValidateBody(c, &mailServer) {
		return
	}

	err := mc.notificationChannelServicer.CheckNotificationChannelConnectivity(context.Background(), mailServer.ToModel())
	if err != nil {
		ginEx.AddError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
