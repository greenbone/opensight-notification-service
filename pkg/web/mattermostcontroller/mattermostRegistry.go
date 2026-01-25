package mattermostcontroller

import (
	"net/http"

	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	svc "github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
)

func ConfigureMappings(r *errmap.Registry) {
	r.Register(
		svc.ErrMattermostChannelLimitReached,
		http.StatusUnprocessableEntity,
		errorResponses.NewErrorGenericResponse("Mattermost channel limit reached."),
	)
	r.Register(
		svc.ErrListMattermostChannels,
		http.StatusInternalServerError,
		errorResponses.ErrorInternalResponse,
	)
}
