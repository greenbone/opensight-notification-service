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

func (m *MattermostNotificationChannelRequest) Cleanup() {
	m.ChannelName = strings.TrimSpace(m.ChannelName)
	m.WebhookUrl = strings.TrimSpace(m.WebhookUrl)
	m.Description = strings.TrimSpace(m.Description)
}

func (m MattermostNotificationChannelRequest) Validate() models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if m.ChannelName == "" {
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

func (r *MattermostNotificationChannelCheckRequest) Cleanup() {
	r.WebhookUrl = strings.TrimSpace(r.WebhookUrl)
}

func (r *MattermostNotificationChannelCheckRequest) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if r.WebhookUrl == "" {
		errs["webhookUrl"] = "webhook URL is required"
	} else {
		_, err := url.Parse(r.WebhookUrl)
		if err != nil && !strings.Contains(err.Error(), "empty url") {
			log.Info().Err(err).Msg("invalid webhook url")
			errs["webhookUrl"] = "invalid"
		}
	}

	return errs
}
