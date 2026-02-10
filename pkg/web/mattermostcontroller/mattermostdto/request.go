package mattermostdto

import (
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/policy"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
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

func (m MattermostNotificationChannelRequest) Validate() *errs.ErrValidation {
	errors := make(map[string]string)

	if m.ChannelName == "" {
		errors["channelName"] = translation.ChannelNameIsRequired
	}

	if m.WebhookUrl == "" {
		errors["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		if _, err := policy.MattermostWebhookUrlPolicy(m.WebhookUrl); err != nil {
			errors["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	if len(errors) > 0 {
		return &errs.ErrValidation{Errors: errors}
	}

	return nil
}

// MattermostNotificationChannelCheckRequest mattermost notification channel check request
type MattermostNotificationChannelCheckRequest struct {
	WebhookUrl string `json:"webhookUrl"`
}

func (r *MattermostNotificationChannelCheckRequest) Cleanup() {
	r.WebhookUrl = strings.TrimSpace(r.WebhookUrl)
}

func (r *MattermostNotificationChannelCheckRequest) Validate() *errs.ErrValidation {
	errors := make(map[string]string)

	if r.WebhookUrl == "" {
		errors["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		if _, err := policy.MattermostWebhookUrlPolicy(r.WebhookUrl); err != nil {
			errors["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	if len(errors) > 0 {
		return &errs.ErrValidation{Errors: errors}
	}

	return nil
}
