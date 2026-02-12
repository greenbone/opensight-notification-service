package mattermostdto

type MattermostNotificationChannelResponse struct {
	Id          string `json:"id"`
	ChannelName string `json:"channelName"`
	WebhookUrl  string `json:"webhookUrl"`
	Description string `json:"description"`
}
