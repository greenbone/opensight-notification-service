package mailcontroller

import (
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckMailServer(t *testing.T) {
	engine := testhelper.NewTestWebEngine()

	notificationChannelServicer := mocks.NewNotificationChannelService(t)

	AddCheckMailServerController(engine, notificationChannelServicer, testhelper.MockAuthMiddlewareWithAdmin)

	t.Run("mail server check is successful", func(t *testing.T) {
		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(nil)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": true,
				"isTlsEnforced": false,
				"username": "testUser",
				"password": "123"
			}`).
			Expect().
			StatusCode(http.StatusNoContent).
			NoContent()
	})

	t.Run("minimal required fields", func(t *testing.T) {
		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"domain": "required",
					"port": "required"
				}
			}`)
	})

	t.Run("username and password are required if isAuthenticationRequired is true", func(t *testing.T) {
		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": true,
				"isTlsEnforced": false,
				"username": "",
				"password": ""
			}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"password": "required",
					"username": "required"
				}
			}`)
	})

	t.Run("return the error if the mail server check fails", func(t *testing.T) {
		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(assert.AnError)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{
				"domain": "example.com",
				"port": 123,
				"isAuthenticationRequired": false,
				"isTlsEnforced": false
			}`).
			Expect().
			StatusCode(http.StatusUnprocessableEntity).
			Json(`{
				"type": "greenbone/generic-error",
				"title": "assert.AnError general error for testing"
			}`)
	})
}
