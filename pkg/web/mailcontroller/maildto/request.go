package maildto

import (
	"net/mail"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
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
		errors["domain"] = "A Mailhub is required."
	}
	if v.Port == 0 {
		errors["port"] = "A port is required."
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errors["username"] = "An Username is required."
		}
		if v.Password == "" {
			errors["password"] = "A Password is required."
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
		errs["domain"] = "required"
	}
	if v.Port == 0 {
		errs["port"] = "required"
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errs["username"] = "required"
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
		errMap["domain"] = "required"
	}

	if r.Port < 1 || r.Port > 65535 {
		errMap["port"] = "required"
	}

	if r.SenderEmailAddress == "" {
		errMap["senderEmailAddress"] = "required"
	}

	_, err := mail.ParseAddress(r.SenderEmailAddress)
	if err != nil && !strings.Contains(err.Error(), "mail: no address") {
		log.Info().Msgf("unable to parse email address %s", err.Error())
		errMap["senderEmailAddress"] = "invalid"
	}

	if r.ChannelName == "" {
		errMap["channelName"] = "required"
	}

	return errMap
}
