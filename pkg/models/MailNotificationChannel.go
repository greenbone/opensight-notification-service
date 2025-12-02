package models

type MailNotificationChannel struct {
	Id                       *string `json:"id,omitempty"`
	ChannelName              *string `json:"channelName" binding:"required"`
	Domain                   *string `json:"domain" binding:"required"`
	Port                     *int    `json:"port" binding:"required"`
	IsAuthenticationRequired *bool   `json:"isAuthenticationRequired" binding:"required"`
	IsTlsEnforced            *bool   `json:"isTlsEnforced" binding:"required"`
	Username                 *string `json:"username" binding:"required"`
	Password                 *string `json:"password" binding:"required"`
	MaxEmailAttachmentSizeMb *int    `json:"maxEmailAttachmentSizeMb" binding:"required"`
	MaxEmailIncludeSizeMb    *int    `json:"maxEmailIncludeSizeMb" binding:"required"`
	SenderEmailAddress       *string `json:"senderEmailAddress" binding:"required,email"`
}
