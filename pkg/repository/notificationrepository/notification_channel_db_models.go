package notificationrepository

import "github.com/greenbone/opensight-notification-service/pkg/models"

type notificationChannelRow struct {
	Id                       string  `db:"id"`
	CreatedAt                string  `db:"created_at"`
	Updated                  *string `db:"updated"`
	ChannelType              string  `db:"channel_type"`
	ChannelName              *string `db:"channel_name"`
	WebhookUrl               *string `db:"webhook_url"`
	Domain                   *string `db:"domain"`
	Port                     *int    `db:"port"`
	IsAuthenticationRequired *bool   `db:"is_authentication_required"`
	IsTlsEnforced            *bool   `db:"is_tls_enforced"`
	Username                 *string `db:"username"`
	Password                 *string `db:"password"`
	MaxEmailAttachmentSizeMb *int    `db:"max_email_attachment_size_mb"`
	MaxEmailIncludeSizeMb    *int    `db:"max_email_include_size_mb"`
	SenderEmailAddress       *string `db:"sender_email_address"`
}

func (r notificationChannelRow) ToModel() models.NotificationChannel {
	// Map fields to your models.NotificationChannel struct
	return models.NotificationChannel{
		Id:                       r.Id,
		CreatedAt:                r.CreatedAt,
		Updated:                  r.Updated,
		ChannelType:              r.ChannelType,
		ChannelName:              r.ChannelName,
		WebhookUrl:               r.WebhookUrl,
		Domain:                   r.Domain,
		Port:                     r.Port,
		IsAuthenticationRequired: r.IsAuthenticationRequired,
		IsTlsEnforced:            r.IsTlsEnforced,
		Username:                 r.Username,
		Password:                 r.Password,
		MaxEmailAttachmentSizeMb: r.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    r.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       r.SenderEmailAddress,
	}
}

// Helper function to map model to DB row struct
func toNotificationChannelRow(in models.NotificationChannel) (notificationChannelRow, error) {
	// Add validation or transformation logic if needed
	return notificationChannelRow{
		Id:                       in.Id,
		CreatedAt:                in.CreatedAt,
		Updated:                  in.Updated,
		ChannelType:              in.ChannelType,
		ChannelName:              in.ChannelName,
		WebhookUrl:               in.WebhookUrl,
		Domain:                   in.Domain,
		Port:                     in.Port,
		IsAuthenticationRequired: in.IsAuthenticationRequired,
		IsTlsEnforced:            in.IsTlsEnforced,
		Username:                 in.Username,
		Password:                 in.Password,
		MaxEmailAttachmentSizeMb: in.MaxEmailAttachmentSizeMb,
		MaxEmailIncludeSizeMb:    in.MaxEmailIncludeSizeMb,
		SenderEmailAddress:       in.SenderEmailAddress,
	}, nil
}
