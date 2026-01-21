package mailcontroller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
)

func TestCheckMailServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	notificationChannelServicer := mocks.NewNotificationChannelService(t)

	AddCheckMailServerController(engine, notificationChannelServicer, testhelper.MockAuthMiddlewareWithAdmin)

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
			StatusCode(400).
			Json(`{
				"type": "greenbone/validation-error",
				"title": "",
				"errors": {
					"password": "required",
					"username": "required"
				}
		}`)
	})

}
