package teamsdto

import (
	"regexp"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
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
		var re = regexp.MustCompile(`^https://[\w.-]+/webhook/[a-zA-Z0-9]+$`)
		if !re.MatchString(m.WebhookUrl) {
			errs["webhookUrl"] = "Invalid teams webhook URL format."
		}
	}

	return errs
}

// TeamsNotificationChannelCheckRequest teams notification channel check request
type TeamsNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *TeamsNotificationChannelCheckRequest) Validate() helper.ValidateErrors {
	errors := make(helper.ValidateErrors)

	if r.WebhookUrl == "" {
		errors["webhookUrl"] = "webhook URL is required"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
