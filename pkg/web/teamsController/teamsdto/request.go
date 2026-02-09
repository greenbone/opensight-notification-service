package teamsdto

import (
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
)

// TeamsNotificationChannelRequest teams notification channel request
type TeamsNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

func (m *TeamsNotificationChannelRequest) Validate() models.ValidationErrors {
	errs := make(map[string]string)
	if m.ChannelName == "" {
		errs["channelName"] = translation.ChannelNameIsRequired
	}

	if m.WebhookUrl == "" {
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	}

	return errs
}

// TeamsNotificationChannelCheckRequest teams notification channel check request
type TeamsNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *TeamsNotificationChannelCheckRequest) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if r.WebhookUrl == "" {
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	}

	return errs
}
