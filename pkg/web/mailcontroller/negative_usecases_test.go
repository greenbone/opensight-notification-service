package mailcontroller

import (
	"net/http"
	"testing"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MailController_Negative_Cases(t *testing.T) {
	t.Parallel()

	valid := testhelper.GetValidMailNotificationChannel()

	t.Run("Check if Create/List/Update Mail Notification doesnt return password", func(t *testing.T) {
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
			JsonPath("$.id", httpassert.ExtractTo(&mailId)).
			JsonPath("$.password", nil)
		require.NotEmpty(t, mailId)

		// --- List ---
		request.Get("/notification-channel/mail").
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$", httpassert.HasSize(1)).
			JsonPath("$[0].password", nil)

		// --- Update ---
		updated := valid
		updated.Id = &mailId
		newName := "updated"
		updated.ChannelName = &newName
		request.Put("/notification-channel/mail/"+mailId).
			JsonContentObject(updated).
			Expect().
			StatusCode(http.StatusOK).
			JsonPath("$.channelName", newName).
			JsonPath("$.password", nil)
	})

	t.Run("Do not update password in Update Mail if it passed as nil", func(t *testing.T) {
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
			return err == nil && password != ""
		}, 5*time.Second, 100*time.Millisecond)

		// --- Update ---
		updated := valid
		updated.Password = nil
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
			return err == nil && password != ""
		}, 5*time.Second, 100*time.Millisecond)
	})
}
