package mailcontroller

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const encryptionPrefix = "ENCV"

func TestIntegration_MailController_CRUD(t *testing.T) {
	valid := testhelper.GetValidMailNotificationChannel()

	t.Run("Perform all the CRUD operations", func(t *testing.T) {
		router, db := setupTestRouter(t)
		defer func() { _ = db.Close() }()

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
			JsonPath("$.port", float64(25)).
			JsonPath("$.isAuthenticationRequired", true).
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
		updated.ChannelName = newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(10)).
			JsonPath("$.id", mailId).
			JsonPath("$.channelName", newName).
			JsonPath("$.domain", "example.com").
			JsonPath("$.port", float64(25)).
			JsonPath("$.isAuthenticationRequired", true).
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

		{ // Create assertion section
			var username, password string
			err := db.QueryRow(`
			SELECT username, password 
			FROM notification_service.notification_channel 
			WHERE id = $1`, mailId).Scan(&username, &password)
			require.NoError(t, err)

			require.NotEmpty(t, username)
			require.NotEmpty(t, password)

			require.True(t, strings.HasPrefix(username, encryptionPrefix), "username encryption prefix is missing")
			require.True(t, strings.HasPrefix(password, encryptionPrefix), "password encryption prefix is missing")
		}

		// --- Update ---
		updated := valid
		newPassword := "newPassword"
		updated.Password = &newPassword
		newName := "updated"
		updated.ChannelName = newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.channelName", newName)

		{ // Update assertions section
			var username, password string
			err := db.QueryRow(`
			SELECT username, password 
			FROM notification_service.notification_channel 
			WHERE id = $1`, mailId).Scan(&username, &password)
			require.NoError(t, err)

			require.NotEmpty(t, username)
			require.NotEmpty(t, password)

			require.True(t, strings.HasPrefix(username, encryptionPrefix), "username encryption prefix is missing")
			require.True(t, strings.HasPrefix(password, encryptionPrefix), "password encryption prefix is missing")
		}
	})
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	mailService := notificationchannelservice.NewMailService()
	mailSvc := notificationchannelservice.NewMailChannelService(svc, repo, mailService, 1)

	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)

	NewMailController(router, svc, mailSvc, testhelper.MockAuthMiddlewareWithAdmin, registry)

	return router, db
}
