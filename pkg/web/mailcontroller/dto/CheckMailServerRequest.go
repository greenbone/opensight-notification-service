package dto

import (
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/web/helper"
)

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

func (v CheckMailServerRequest) Validate() helper.ValidateErrors {
	errors := make(helper.ValidateErrors)

	if v.Domain == "" {
		errors["domain"] = "required"
	}
	if v.Port == 0 {
		errors["port"] = "required"
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errors["username"] = "required"
		}
		if v.Password == "" {
			errors["password"] = "required"
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

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

func (v CheckMailServerEntityRequest) Validate() helper.ValidateErrors {
	errors := make(helper.ValidateErrors)

	if v.Domain == "" {
		errors["domain"] = "required"
	}
	if v.Port == 0 {
		errors["port"] = "required"
	}

	if v.IsAuthenticationRequired {
		if v.Username == "" {
			errors["username"] = "required"
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

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
