package teamsdto

import (
	"net/url"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
)

// TeamsNotificationChannelRequest teams notification channel request
type TeamsNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

func (m *TeamsNotificationChannelRequest) Cleanup() {
	m.ChannelName = strings.TrimSpace(m.ChannelName)
	m.WebhookUrl = strings.TrimSpace(m.WebhookUrl)
	m.Description = strings.TrimSpace(m.Description)
}

func (m *TeamsNotificationChannelRequest) Validate() models.ValidationErrors {
	errs := make(map[string]string)

	if m.ChannelName == "" {
		errs["channelName"] = translation.ChannelNameIsRequired
	}

	if m.WebhookUrl == "" {
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		_, err := url.ParseRequestURI(m.WebhookUrl)
		if err != nil {
			errs["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	return errs
}

// TeamsNotificationChannelCheckRequest teams notification channel check request
type TeamsNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (m *TeamsNotificationChannelCheckRequest) Cleanup() {
	m.WebhookUrl = strings.TrimSpace(m.WebhookUrl)
}

func (r *TeamsNotificationChannelCheckRequest) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if r.WebhookUrl == "" {
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		_, err := url.ParseRequestURI(r.WebhookUrl)
		if err != nil {
			errs["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	return errs
}
