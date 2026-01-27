package mailcontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup(t *testing.T) (*gin.Engine, *mocks.MailChannelService) {
	engine := testhelper.NewTestWebEngine()

	notificationChannelServicer := mocks.NewMailChannelService(t)

	AddCheckMailServerController(engine, notificationChannelServicer, testhelper.MockAuthMiddlewareWithAdmin)
	return engine, notificationChannelServicer
}

func TestCheckMailServer(t *testing.T) {
	t.Run("mail server check is successful", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

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
		engine, _ := setup(t)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`{}`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"domain": "A Mailhub is required.",
					"port": "A port is required."
				}
			}`)
	})

	t.Run("username and password are required if isAuthenticationRequired is true", func(t *testing.T) {
		engine, _ := setup(t)

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
					"username": "An Username is required.",
					"password": "A Password is required."
				}
			}`)
	})

	t.Run("return the error if the mail server check fails", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

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
