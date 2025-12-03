package models

type MailNotificationChannel struct {
	Id                       *string `json:"id,omitempty"`
	ChannelName              *string `json:"channelName" binding:"required"`
	Domain                   *string `json:"domain" binding:"required"`
	Port                     *int    `json:"port" binding:"required"`
	IsAuthenticationRequired *bool   `json:"isAuthenticationRequired" binding:"required"`
	IsTlsEnforced            *bool   `json:"isTlsEnforced" binding:"required"`
	Username                 *string `json:"username,omitempty"`
	Password                 *string `json:"password,omitempty"`
	MaxEmailAttachmentSizeMb *int    `json:"maxEmailAttachmentSizeMb,omitempty"`
	MaxEmailIncludeSizeMb    *int    `json:"maxEmailIncludeSizeMb,omitempty"`
	SenderEmailAddress       *string `json:"senderEmailAddress" binding:"required,email"`
}
