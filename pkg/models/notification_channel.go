package models

type NotificationChannel struct {
	Id                       string  `json:"id" readonly:"true"`
	CreatedAt                string  `json:"createdAt" readonly:"true"`
	Updated                  *string `json:"updated,omitempty"`
	ChannelType              string  `json:"channelType" binding:"required"`
	ChannelName              *string `json:"channelName,omitempty"`
	WebhookUrl               *string `json:"webhookUrl,omitempty"`
	Domain                   *string `json:"domain,omitempty"`
	Port                     *int    `json:"port,omitempty"`
	IsAuthenticationRequired *bool   `json:"isAuthenticationRequired,omitempty"`
	IsTlsEnforced            *bool   `json:"isTlsEnforced,omitempty"`
	Username                 *string `json:"username,omitempty"`
	Password                 *string `json:"password,omitempty"`
	MaxEmailAttachmentSizeMb *int    `json:"maxEmailAttachmentSizeMb,omitempty"`
	MaxEmailIncludeSizeMb    *int    `json:"maxEmailIncludeSizeMb,omitempty"`
	SenderEmailAddress       *string `json:"senderEmailAddress,omitempty"`
}
