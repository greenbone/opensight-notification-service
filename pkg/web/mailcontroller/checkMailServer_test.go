package mailcontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
)

func setup(t *testing.T) (*gin.Engine, *mocks.MailChannelService) {
	registry := errmap.NewRegistry()
	engine := testhelper.NewTestWebEngine(registry)

	notificationChannelServicer := mocks.NewMailChannelService(t)

	AddCheckMailServerController(engine, notificationChannelServicer, testhelper.MockAuthMiddlewareWithAdmin, registry)
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

	t.Run("none request body", func(t *testing.T) {
		engine, _ := setup(t)

		httpassert.New(t, engine).
			Post("/notification-channel/mail/check").
			Content(`-`).
			Expect().
			StatusCode(http.StatusBadRequest).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "unable to parse the request",
				"details":"error parsing body"
			}`)
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
					"domain": "A mailhub is required.",
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
					"username": "A username is required.",
					"password": "A password is required."
				}
			}`)
	})

	t.Run("return the error if the mail server check fails", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(notificationchannelservice.ErrMailServerUnreachable)

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
				"title": "Server is unreachable"
			}`)
	})
}
