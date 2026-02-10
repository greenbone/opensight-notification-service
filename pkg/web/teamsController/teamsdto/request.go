package teamsdto

import (
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
)

// TeamsNotificationChannelRequest teams notification channel request
type TeamsNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}

func (m *TeamsNotificationChannelRequest) Validate() *errs.ErrValidation {
	errors := make(map[string]string)
	if m.ChannelName == "" {
		errors["channelName"] = translation.ChannelNameIsRequired
	}

	if m.WebhookUrl == "" {
		errors["webhookUrl"] = translation.WebhookUrlIsRequired
	}

	if len(errors) > 0 {
		return &errs.ErrValidation{Errors: errors}
	}

	return nil
}

// TeamsNotificationChannelCheckRequest teams notification channel check request
type TeamsNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *TeamsNotificationChannelCheckRequest) Validate() *errs.ErrValidation {
	errors := make(map[string]string)

	if r.WebhookUrl == "" {
		errors["webhookUrl"] = translation.WebhookUrlIsRequired
	}

	if len(errors) > 0 {
		return &errs.ErrValidation{Errors: errors}
	}

	return nil
}
