package mapper

import (
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMapNotificationChannelToMail(t *testing.T) {
	channel := models.NotificationChannel{
		Id:                       ptrString("id1"),
		ChannelName:              ptrString("TestChannel"),
		Domain:                   ptrString("example.com"),
		Port:                     ptrInt(587),
		IsAuthenticationRequired: ptrBool(true),
		IsTlsEnforced:            ptrBool(true),
		Username:                 ptrString("user"),
		MaxEmailAttachmentSizeMb: ptrInt(10),
		MaxEmailIncludeSizeMb:    ptrInt(5),
		SenderEmailAddress:       ptrString("sender@example.com"),
	}

	mail := MapNotificationChannelToMail(channel)

	t.Run("assert all fields", func(t *testing.T) {
		assert.True(t, mail.Id == channel.Id)
		assert.True(t, mail.ChannelName == channel.ChannelName)
		assert.True(t, mail.Domain == channel.Domain)
		assert.True(t, mail.Port == channel.Port)
		assert.True(t, mail.IsAuthenticationRequired == channel.IsAuthenticationRequired)
		assert.True(t, mail.IsTlsEnforced == channel.IsTlsEnforced)
		assert.True(t, mail.Username == channel.Username)
		assert.True(t, mail.MaxEmailAttachmentSizeMb == channel.MaxEmailAttachmentSizeMb)
		assert.True(t, mail.MaxEmailIncludeSizeMb == channel.MaxEmailIncludeSizeMb)
		assert.True(t, mail.SenderEmailAddress == channel.SenderEmailAddress)
	})
}

func TestMapMailToNotificationChannel(t *testing.T) {
	mail := models.MailNotificationChannel{
		Id:                       ptrString("id2"),
		ChannelName:              ptrString("MailChannel"),
		Domain:                   ptrString("mail.com"),
		Port:                     ptrInt(465),
		IsAuthenticationRequired: ptrBool(false),
		IsTlsEnforced:            ptrBool(false),
		Username:                 ptrString("mailuser"),
		Password:                 ptrString("secret"),
		MaxEmailAttachmentSizeMb: ptrInt(20),
		MaxEmailIncludeSizeMb:    ptrInt(15),
		SenderEmailAddress:       ptrString("mail@domain.com"),
	}

	channel := MapMailToNotificationChannel(mail)

	t.Run("assert all fields", func(t *testing.T) {
		assert.True(t, channel.ChannelType == "mail")
		assert.True(t, channel.Id == mail.Id)
		assert.True(t, channel.ChannelName == mail.ChannelName)
		assert.True(t, channel.Domain == mail.Domain)
		assert.True(t, channel.Port == mail.Port)
		assert.True(t, channel.IsAuthenticationRequired == mail.IsAuthenticationRequired)
		assert.True(t, channel.IsTlsEnforced == mail.IsTlsEnforced)
		assert.True(t, channel.Username == mail.Username)
		assert.True(t, channel.Password == mail.Password)
		assert.True(t, channel.MaxEmailAttachmentSizeMb == mail.MaxEmailAttachmentSizeMb)
		assert.True(t, channel.MaxEmailIncludeSizeMb == mail.MaxEmailIncludeSizeMb)
		assert.True(t, channel.SenderEmailAddress == mail.SenderEmailAddress)
	})
}

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrBool(b bool) *bool       { return &b }
