package mailcontroller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestIntegration_MailController_CRUD(t *testing.T) {
	// Use testcontainer-based DB
	gormDB := testhelper.NewTestDB(t)
	sqlxdb, err := gormDB.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlxdb.Close() // Ensure DB is closed on exit

	db := sqlx.NewDb(sqlxdb, "postgres")
	db.Exec("DELETE FROM notification_service.notification_channel")

	repo, err := notificationrepository.NewNotificationChannelRepository(db)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	svc := notificationchannelservice.NewNotificationChannelService(repo)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewMailController(router, svc, testhelper.MockAuthMiddlewareWithAdmin)
	valid := getValidNotificationChannel()

	// --- Create ---
	createResp := httpassert.New(t, router).
		Post("/notification-channel/mail").
		JsonContentObject(valid).
		Expect()
	createResp.StatusCode(http.StatusCreated)
	var created models.MailNotificationChannel

	body := createResp.GetBody()
	if err := json.Unmarshal([]byte(body), &created); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	if created.Id == nil || created.ChannelName == nil || *created.ChannelName != *valid.ChannelName {
		t.Fatalf("unexpected create response: %+v", created)
	}
	valid.Id = created.Id

	// --- List ---
	listResp := httpassert.New(t, router).
		Get("/notification-channel/mail").
		Expect()
	listResp.StatusCode(http.StatusOK)
	listResp.JsonPath("$", httpassert.HasSize(1))

	// --- Update ---
	updated := valid
	newName := "updated"
	updated.ChannelName = &newName
	updateResp := httpassert.New(t, router).
		Put("/notification-channel/mail/" + *created.Id).
		JsonContentObject(updated).
		Expect()
	updateResp.StatusCode(http.StatusOK)
	updateResp.JsonPath("$.channelName", newName)

	// --- Delete ---
	delResp := httpassert.New(t, router).
		Delete("/notification-channel/mail/" + *created.Id).
		Expect()
	delResp.StatusCode(http.StatusNoContent)

	// --- List after delete ---
	listAfterDel := httpassert.New(t, router).
		Get("/notification-channel/mail").
		Expect()
	listAfterDel.StatusCode(http.StatusOK)
	listAfterDel.JsonPath("$", httpassert.HasSize(0))
}
