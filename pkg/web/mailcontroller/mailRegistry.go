package mailcontroller

import (
	"net/http"

	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	svc "github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
)

func ConfigureMappings(r *errmap.Registry) {
	r.Register(
		svc.ErrMailChannelLimitReached,
		http.StatusConflict, // TODO check agreement should we use `StatusUnprocessableEntity` instead??
		errorResponses.NewErrorGenericResponse("Mail channel limit reached."),
	)

	r.Register(
		svc.ErrCreateMailFailed,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Unable to create mail client"),
	)

	r.Register(
		svc.ErrMailServerUnreachable,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Server is unreachable"),
	)

	r.Register(
		svc.ErrGetMailChannel,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)

	r.Register(
		svc.ErrListMailChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
}
