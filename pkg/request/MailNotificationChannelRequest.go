package request

import (
	"net/mail"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/rs/zerolog/log"
)

type MailNotificationChannelRequest struct {
	Id                       *string `json:"id,omitempty"`
	ChannelName              string  `json:"channelName"`
	Domain                   string  `json:"domain"`
	Port                     int     `json:"port"`
	IsAuthenticationRequired bool    `json:"isAuthenticationRequired" default:"false"`
	IsTlsEnforced            bool    `json:"isTlsEnforced" default:"false"`
	Username                 *string `json:"username,omitempty"`
	Password                 *string `json:"password,omitempty"`
	MaxEmailAttachmentSizeMb *int    `json:"maxEmailAttachmentSizeMb,omitempty"`
	MaxEmailIncludeSizeMb    *int    `json:"maxEmailIncludeSizeMb,omitempty"`
	SenderEmailAddress       string  `json:"senderEmailAddress"`
}

func (r MailNotificationChannelRequest) WithEmptyPassword() MailNotificationChannelRequest {
	r.Password = nil
	return r
}

func (r MailNotificationChannelRequest) Validate() models.ValidationErrors {
	errMap := make(models.ValidationErrors)

	if strings.TrimSpace(r.Domain) == "" {
		errMap["domain"] = "required"
	}

	if r.Port < 1 || r.Port > 65535 {
		errMap["port"] = "required"
	}

	if strings.TrimSpace(r.SenderEmailAddress) == "" {
		errMap["senderEmailAddress"] = "required"
	}

	_, err := mail.ParseAddress(r.SenderEmailAddress)
	if err != nil && !strings.Contains(err.Error(), "mail: no address") {
		log.Info().Msgf("unable to parse email address %s", err.Error())
		errMap["senderEmailAddress"] = "invalid"
	}

	if strings.TrimSpace(r.ChannelName) == "" {
		errMap["channelName"] = "required"
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}
