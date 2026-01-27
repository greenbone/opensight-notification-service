package teamsdto

import "github.com/greenbone/opensight-notification-service/pkg/web/helper"

type TeamsNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

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
