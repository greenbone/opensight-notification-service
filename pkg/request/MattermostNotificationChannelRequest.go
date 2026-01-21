package request

type MattermostNotificationChannelRequest struct {
	Id          *string `json:"id,omitempty"`
	ChannelName string  `json:"channelName"`
	WebhookUrl  string  `json:"webhookUrl"`
	Description string  `json:"description"`
}
