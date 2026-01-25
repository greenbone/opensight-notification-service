package request

import (
	"net/url"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/rs/zerolog/log"
)

type MattermostNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

func (m *MattermostNotificationChannelRequest) Validate() models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if strings.TrimSpace(m.ChannelName) == "" {
		errors["webhookUrl"] = "required"
	}

	_, err := url.Parse(m.WebhookUrl)
	if err != nil && !strings.Contains(err.Error(), "empty url") {
		log.Info().Err(err).Msg("invalid webhook url")
		errors["webhookUrl"] = "invalid"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
