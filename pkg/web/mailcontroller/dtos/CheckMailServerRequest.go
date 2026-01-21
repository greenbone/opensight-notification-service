package dtos

import "github.com/greenbone/opensight-notification-service/pkg/models"

type CheckMailServerRequest struct {
	Domain                   string `json:"domain"`
	Port                     int    `json:"port"`
	IsAuthenticationRequired bool   `json:"isAuthenticationRequired" default:"false"`
	IsTlsEnforced            bool   `json:"isTlsEnforced" default:"false"`
	Username                 string `json:"username"`
	Password                 string `json:"password"`
}

func NewCheckMailServerRequest(channel models.NotificationChannel) CheckMailServerRequest {
	return CheckMailServerRequest{
		Domain:                   *channel.Domain,
		Port:                     *channel.Port,
		IsAuthenticationRequired: *channel.IsAuthenticationRequired,
		IsTlsEnforced:            *channel.IsTlsEnforced,
		Username:                 *channel.Username,
		Password:                 *channel.Password,
	}
}

type ValidateErrors map[string]string

func (v ValidateErrors) Error() string {
	return "validation error"
}

func (v CheckMailServerRequest) Validate() ValidateErrors {
	errors := make(ValidateErrors)

	if v.Domain == "" {
		errors["domain"] = "required"
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
