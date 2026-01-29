package maildto

import (
	"net/mail"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/translation"
	"github.com/rs/zerolog/log"
)

// CheckMailServerRequest check mail server request
type CheckMailServerRequest struct {
	Domain                   string `json:"domain"`
	Port                     int    `json:"port"`
	IsAuthenticationRequired bool   `json:"isAuthenticationRequired" default:"false"`
	IsTlsEnforced            bool   `json:"isTlsEnforced" default:"false"`
	Username                 string `json:"username"`
	Password                 string `json:"password"`
}

func (v CheckMailServerRequest) ToModel() models.NotificationChannel {
	return models.NotificationChannel{
		Domain:                   &v.Domain,
		Port:                     &v.Port,
		IsAuthenticationRequired: &v.IsAuthenticationRequired,
		IsTlsEnforced:            &v.IsTlsEnforced,
		Username:                 &v.Username,
		Password:                 &v.Password,
	}
}

func (r *CheckMailServerRequest) Cleanup() {
	r.Domain = strings.TrimSpace(r.Domain)
}

func (v CheckMailServerRequest) Validate() models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if v.Domain == "" {
		errors["domain"] = translation.MailhubIsRequired
	}
	if v.Port == 0 {
		errors["port"] = translation.PortIsRequired
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errors["username"] = translation.UsernameIsRequired
		}
		if v.Password == "" {
			errors["password"] = translation.PasswordIsRequired
		}
	}

	return errors
}

// CheckMailServerEntityRequest check mail server entity request
type CheckMailServerEntityRequest struct {
	Domain                   string `json:"domain"`
	Port                     int    `json:"port"`
	IsAuthenticationRequired bool   `json:"isAuthenticationRequired" default:"false"`
	IsTlsEnforced            bool   `json:"isTlsEnforced" default:"false"`
	Username                 string `json:"username"`
	Password                 string `json:"password"`
}

func (v CheckMailServerEntityRequest) ToModel() models.NotificationChannel {
	return models.NotificationChannel{
		Domain:                   &v.Domain,
		Port:                     &v.Port,
		IsAuthenticationRequired: &v.IsAuthenticationRequired,
		IsTlsEnforced:            &v.IsTlsEnforced,
		Username:                 &v.Username,
		Password:                 &v.Password,
	}
}

func (r *CheckMailServerEntityRequest) Cleanup() {
	r.Domain = strings.TrimSpace(r.Domain)
}

func (v CheckMailServerEntityRequest) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if v.Domain == "" {
		errs["domain"] = translation.MailhubIsRequired
	}
	if v.Port < 1 || v.Port > 65535 {
		errs["port"] = translation.PortIsRequired
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errs["username"] = translation.UsernameIsRequired
		}
	}

	return errs
}

// MailNotificationChannelRequest mail notification channel request
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

func (r *MailNotificationChannelRequest) Cleanup() {
	r.Domain = strings.TrimSpace(r.Domain)
	r.SenderEmailAddress = strings.TrimSpace(r.SenderEmailAddress)
	r.ChannelName = strings.TrimSpace(r.ChannelName)
}

func (r MailNotificationChannelRequest) Validate() models.ValidationErrors {
	errMap := make(models.ValidationErrors)

	if r.Domain == "" {
		errMap["domain"] = translation.MailhubIsRequired
	}

	if r.Port < 1 || r.Port > 65535 {
		errMap["port"] = translation.PortIsRequired
	}

	if r.SenderEmailAddress == "" {
		errMap["senderEmailAddress"] = translation.MailSenderIsRequired
	}

	_, err := mail.ParseAddress(r.SenderEmailAddress)
	if err != nil && !strings.Contains(err.Error(), "mail: no address") {
		log.Info().Msgf("unable to parse email address %s", err.Error())
		errMap["senderEmailAddress"] = translation.ValidEmailSenderIsRequired
	}

	if r.ChannelName == "" {
		errMap["channelName"] = translation.ChannelNameIsRequired
	}

	return errMap
}
