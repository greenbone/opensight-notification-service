package mattermostdto

import (
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
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

func (m MattermostNotificationChannelRequest) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if m.ChannelName == "" {
		errs["channelName"] = translation.ChannelNameIsRequired
	}

	if m.WebhookUrl == "" {
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		if _, err := policy.MattermostWebhookUrlPolicy(m.WebhookUrl); err != nil {
			errs["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	return errs
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
		errs["webhookUrl"] = translation.WebhookUrlIsRequired
	} else {
		if _, err := policy.MattermostWebhookUrlPolicy(r.WebhookUrl); err != nil {
			errs["webhookUrl"] = translation.ValidWebhookUrlIsRequired
		}
	}

	return errs
}
