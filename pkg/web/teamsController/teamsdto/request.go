package teamsdto

import (
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/policy"
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
		errs["channelName"] = "A channel name is required."
	}

	if m.WebhookUrl == "" {
		errs["webhookUrl"] = "A Webhook URL is required."
	} else {
		if _, err := policy.TeamsWebhookUrlPolicy(m.WebhookUrl); err != nil {
			errs["webhookUrl"] = "Please enter a valid webhook URL."
		}
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
		errs["webhookUrl"] = "A Webhook URL is required."
	} else {
		if _, err := policy.TeamsWebhookUrlPolicy(r.WebhookUrl); err != nil {
			errs["webhookUrl"] = "Please enter a valid webhook URL."
		}
	}

	return errs
}
