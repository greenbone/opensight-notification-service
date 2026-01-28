package maildto

type MailNotificationChannelResponse struct {
	Id                       string  `json:"id,omitempty"`
	ChannelName              string  `json:"channelName"`
	Domain                   string  `json:"domain"`
	Port                     int     `json:"port"`
	IsAuthenticationRequired bool    `json:"isAuthenticationRequired" default:"false"`
	IsTlsEnforced            bool    `json:"isTlsEnforced" default:"false"`
	Username                 *string `json:"username,omitempty"`
	MaxEmailAttachmentSizeMb *int    `json:"maxEmailAttachmentSizeMb,omitempty"`
	MaxEmailIncludeSizeMb    *int    `json:"maxEmailIncludeSizeMb,omitempty"`
	SenderEmailAddress       string  `json:"senderEmailAddress"`
}
