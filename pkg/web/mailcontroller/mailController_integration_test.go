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
		t.Parallel()
		router, db := setupTestRouter(t)
		defer db.Close()
		request := httpassert.New(t, router)

		// --- Create ---
		var mailId string
		request.Post("/notification-channel/mail").
			JsonContentObject(valid).
			Expect().
			StatusCode(http.StatusCreated).
			JsonPath("$.id", httpassert.ExtractTo(&mailId))
		require.NotEmpty(t, mailId)

		// --- List ---
		request.Get("/notification-channel/mail").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1))

		// --- Update ---
		updated := valid
		updated.Id = &mailId
		newName := "updated"
		updated.ChannelName = &newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.channelName", newName)

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
		t.Parallel()
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
			err := db.QueryRow("SELECT password FROM notification_service.notification_channel WHERE id = $1", mailId).Scan(&password)
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
			err := db.QueryRow("SELECT password FROM notification_service.notification_channel WHERE id = $1", mailId).Scan(&password)
			return err == nil && password == newPassword
		}, 5*time.Second, 100*time.Millisecond)
	})
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sqlx.DB) {
	repo, db := testhelper.SetupNotificationChannelTestEnv(t)
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewMailController(router, svc, testhelper.MockAuthMiddlewareWithAdmin)

	return router, db
}
