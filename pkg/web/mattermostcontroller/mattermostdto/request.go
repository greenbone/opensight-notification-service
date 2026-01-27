package mattermostdto

import (
	"net/url"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/rs/zerolog/log"
)

// MattermostNotificationChannelRequest mattermost notification channel request
type MattermostNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

func (m MattermostNotificationChannelRequest) Validate() models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if strings.TrimSpace(m.ChannelName) == "" {
		errors["webhookUrl"] = "required"
	}

	_, err := url.Parse(m.WebhookUrl)
	if err != nil && !strings.Contains(err.Error(), "empty url") {
		log.Info().Err(err).Msg("invalid webhook url")
		errors["webhookUrl"] = "invalid"
	}

	return errors
}

// MattermostNotificationChannelCheckRequest mattermost notification channel check request
type MattermostNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *MattermostNotificationChannelCheckRequest) Validate() models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if r.WebhookUrl == "" {
		errors["webhookUrl"] = "webhook URL is required"
	}

	return errors
}
