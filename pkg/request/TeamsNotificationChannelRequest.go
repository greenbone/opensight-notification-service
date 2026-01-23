package request

type TeamsNotificationChannelRequest struct {
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}
