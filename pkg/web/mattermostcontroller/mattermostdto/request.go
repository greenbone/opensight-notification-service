package mattermostdto

import "github.com/greenbone/opensight-notification-service/pkg/web/helper"

type MattermostNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

type MattermostNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *MattermostNotificationChannelCheckRequest) Validate() helper.ValidateErrors {
	errors := make(helper.ValidateErrors)

	if r.WebhookUrl == "" {
		errors["webhookUrl"] = "webhook URL is required"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
