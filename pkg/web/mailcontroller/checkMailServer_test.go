package mailcontroller

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	svc "github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/middleware"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
)

func setup(t *testing.T) (*gin.Engine, *mocks.NotificationChannelService) {
	registry := errmap.NewRegistry()

	engine := testhelper.NewTestWebEngine()
	engine.Use(middleware.InterpretErrors(gin.ErrorTypePrivate, registry))

	notificationChannelServicer := mocks.NewNotificationChannelService(t)
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
				"details":"unexpected EOF"
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
					"domain": "required",
					"port": "required"
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
					"password": "required",
					"username": "required"
				}
			}`)
	})

	t.Run("return the error if the mail server check fails", func(t *testing.T) {
		engine, notificationChannelServicer := setup(t)

		notificationChannelServicer.EXPECT().CheckNotificationChannelConnectivity(mock.Anything, mock.Anything).
			Return(svc.ErrMailServerUnreachable)

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
