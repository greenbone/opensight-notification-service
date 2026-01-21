//go:build integration
// +build integration

package mailcontroller

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MailController_CRUD(t *testing.T) {
	t.Parallel()

	valid := testhelper.GetValidMailNotificationChannel()

	t.Run("Perform all the CRUD operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		// --- Create ---
		var mailId string
		request.Post("/notification-channel/mail").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$", httpassert.HasSize(10)).
			JsonPath("$.id", httpassert.ExtractTo(&mailId)).
			JsonPath("$.channelName", "mail1").
			JsonPath("$.domain", "example.com").
			JsonPath("$.port", float64(25)).JsonPath("$.isAuthenticationRequired", true).
			JsonPath("$.isTlsEnforced", false).
			JsonPath("$.username", "user").
			JsonPath("$.maxEmailAttachmentSizeMb", float64(10)).
			JsonPath("$.maxEmailIncludeSizeMb", float64(5)).
			JsonPath("$.senderEmailAddress", "sender@example.com")
		require.NotEmpty(t, mailId)

		// --- List ---
		request.Get("/notification-channel/mail").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1)).
			JsonPath("$[0]", httpassert.HasSize(10)).
			JsonPath("$[0].id", httpassert.ExtractTo(&mailId)).
			JsonPath("$[0].channelName", "mail1").
			JsonPath("$[0].domain", "example.com").
			JsonPath("$[0].port", float64(25)).
			JsonPath("$[0].isAuthenticationRequired", true).
			JsonPath("$[0].isTlsEnforced", false).
			JsonPath("$[0].username", "user").
			JsonPath("$[0].maxEmailAttachmentSizeMb", float64(10)).
			JsonPath("$[0].maxEmailIncludeSizeMb", float64(5)).
			JsonPath("$[0].senderEmailAddress", "sender@example.com")

		// --- Update ---
		updated := valid
		updated.Id = &mailId
		newName := "updated"
		updated.ChannelName = &newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(10)).
			JsonPath("$.id", mailId).
			JsonPath("$.channelName", newName).
			JsonPath("$.domain", "example.com").
			JsonPath("$.port", float64(25)).JsonPath("$.isAuthenticationRequired", true).
			JsonPath("$.isTlsEnforced", false).
			JsonPath("$.username", "user").
			JsonPath("$.maxEmailAttachmentSizeMb", float64(10)).
			JsonPath("$.maxEmailIncludeSizeMb", float64(5)).
			JsonPath("$.senderEmailAddress", "sender@example.com")

		// --- Delete ---
		request.Delete("/notification-channel/mail/" + mailId).
			Expect().
			StatusCode(http.StatusNoContent)

		// --- List after delete ---
		request.Get("/notification-channel/mail").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(0))
	})

	t.Run("Update password with Update Mail", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer db.Close()

		request := httpassert.New(t, router)
		var mailId string

		// --- Create ---
		request.Post("/notification-channel/mail").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mailId))
		require.NotEmpty(t, mailId)

		require.Eventually(t, func() bool {
			var password string
			err := db.QueryRow("SELECT password FROM notification_service.notification_channel WHERE id = $1",
				mailId).Scan(&password)
			return err == nil && password == "pass"
		}, 5*time.Second, 100*time.Millisecond)

		// --- Update ---
		updated := valid
		newPassword := "newPassword"
		updated.Password = &newPassword
		newName := "updated"
		updated.ChannelName = &newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.channelName", newName)

		require.Eventually(t, func() bool {
			var password string
			err := db.QueryRow("SELECT password FROM notification_service.notification_channel WHERE id = $1",
				mailId).Scan(&password)
			return err == nil && password == newPassword
		}, 5*time.Second, 100*time.Millisecond)
	})
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	mailSvc := notificationchannelservice.NewMailChannelService(svc)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewMailController(router, svc, mailSvc, testhelper.MockAuthMiddlewareWithAdmin)

	return router, db
}
