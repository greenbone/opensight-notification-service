package mapper

import (
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMapNotificationChannelToMail(t *testing.T) {
	channel := models.NotificationChannel{
		Id:                       helper.ToPtr("id1"),
		ChannelName:              helper.ToPtr("TestChannel"),
		Domain:                   helper.ToPtr("example.com"),
		Port:                     helper.ToPtr("587"),
		IsAuthenticationRequired: helper.ToPtr(true),
		IsTlsEnforced:            helper.ToPtr(true),
		Username:                 helper.ToPtr("user"),
		MaxEmailAttachmentSizeMb: helper.ToPtr(10),
		MaxEmailIncludeSizeMb:    helper.ToPtr(5),
		SenderEmailAddress:       helper.ToPtr("sender@example.com"),
	}

	mail := MapNotificationChannelToMail(channel)

	t.Run("assert all fields", func(t *testing.T) {
		assert.Equal(t, channel.Id, mail.Id)
		assert.Equal(t, channel.ChannelName, &mail.ChannelName)
		assert.Equal(t, channel.Domain, &mail.Domain)
		assert.Equal(t, channel.Port, &mail.Port)
		assert.Equal(t, channel.IsAuthenticationRequired, &mail.IsAuthenticationRequired)
		assert.Equal(t, channel.IsTlsEnforced, &mail.IsTlsEnforced)
		assert.Equal(t, channel.Username, mail.Username)
		assert.Equal(t, channel.MaxEmailAttachmentSizeMb, mail.MaxEmailAttachmentSizeMb)
		assert.Equal(t, channel.MaxEmailIncludeSizeMb, mail.MaxEmailIncludeSizeMb)
		assert.Equal(t, channel.SenderEmailAddress, &mail.SenderEmailAddress)
	})
}

func TestMapMailToNotificationChannel(t *testing.T) {
	mail := models.MailNotificationChannel{
		Id:                       helper.ToPtr("id2"),
		ChannelName:              "MailChannel",
		Domain:                   "mail.com",
		Port:                     "465",
		IsAuthenticationRequired: false,
		IsTlsEnforced:            false,
		Username:                 helper.ToPtr("mailuser"),
		Password:                 helper.ToPtr("secret"),
		MaxEmailAttachmentSizeMb: helper.ToPtr(20),
		MaxEmailIncludeSizeMb:    helper.ToPtr(15),
		SenderEmailAddress:       "mail@domain.com",
	}

	channel := MapMailToNotificationChannel(mail)

	t.Run("assert all fields", func(t *testing.T) {
		assert.Equal(t, "mail", channel.ChannelType)
		assert.Equal(t, mail.Id, channel.Id)
		assert.Equal(t, mail.ChannelName, *channel.ChannelName)
		assert.Equal(t, mail.Domain, *channel.Domain)
		assert.Equal(t, mail.Port, *channel.Port)
		assert.Equal(t, mail.IsAuthenticationRequired, *channel.IsAuthenticationRequired)
		assert.Equal(t, mail.IsTlsEnforced, *channel.IsTlsEnforced)
		assert.Equal(t, mail.Username, channel.Username)
		assert.Equal(t, mail.Password, channel.Password)
		assert.Equal(t, mail.MaxEmailAttachmentSizeMb, channel.MaxEmailAttachmentSizeMb)
		assert.Equal(t, mail.MaxEmailIncludeSizeMb, channel.MaxEmailIncludeSizeMb)
		assert.Equal(t, mail.SenderEmailAddress, *channel.SenderEmailAddress)
	})
}
