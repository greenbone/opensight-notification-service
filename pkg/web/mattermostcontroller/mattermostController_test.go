package mattermostcontroller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/mattermostdto"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

func setupTestController() (*gin.Engine, *mocks.NotificationChannelService, *mocks.MattermostChannelService) {
	r := testhelper.NewTestWebEngine()

	notificationService := &mocks.NotificationChannelService{}
	mattermostService := &mocks.MattermostChannelService{}
	NewMattermostController(r, notificationService, mattermostService, testhelper.MockAuthMiddlewareWithAdmin)

	return r, notificationService, mattermostService
}

func TestCreateMattermostChannel_Success(t *testing.T) {
	r, _, mattermostService := setupTestController()
	input := mattermostdto.MattermostNotificationChannelRequest{
		ChannelName: "test-channel",
		WebhookUrl:  "https://webhookurl.com/hooks/id1",
		Description: "desc",
	}
	output := mattermostdto.MattermostNotificationChannelResponse{
		ChannelName: "test-channel",
		WebhookUrl:  "https://webhookurl.com/hooks/id1",
		Description: "desc",
		Id:          "1",
	}
	mattermostService.On("CreateMattermostChannel", mock.Anything, input).Return(output, nil)

	httpassert.New(t, r).
		Post("/notification-channel/mattermost").
		JsonContentObject(input).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$.channelName", input.ChannelName)
}

func TestCreateMattermostChannel_BadRequest(t *testing.T) {
	r, _, _ := setupTestController()
	badBody := []byte(`{"invalid":}`)
	httpassert.New(t, r).
		Post("/notification-channel/mattermost").
		JsonContentObject(badBody).
		Expect().
		StatusCode(http.StatusBadRequest)
}

func TestListMattermostChannelsByType_Success(t *testing.T) {
	r, notificationService, _ := setupTestController()
	id := helper.ToPtr("123")
	desc := helper.ToPtr("desc")
	channels := []models.NotificationChannel{{
		Id:          id,
		ChannelName: helper.ToPtr("test"),
		WebhookUrl:  helper.ToPtr("url"),
		Description: desc,
	}}
	notificationService.
		On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).
		Return(channels, nil)

	httpassert.New(t, r).
		Get("/notification-channel/mattermost").
		Expect().
		StatusCode(http.StatusOK).
		JsonPath("$", httpassert.HasSize(1)).
		JsonPath("$[0].channelName", "test").
		JsonPath("$[0].webhookUrl", "url").
		JsonPath("$[0].description", "desc").
		JsonPath("$[0].id", "123")
}

func TestListMattermostChannelsByType_Error(t *testing.T) {
	r, notificationService, _ := setupTestController()
	notificationService.
		On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).
		Return(nil, errors.New("fail"))

	httpassert.New(t, r).
		Get("/notification-channel/mattermost").
		Expect().
		StatusCode(http.StatusInternalServerError)
}

func TestUpdateMattermostChannel_Success(t *testing.T) {
	r, notificationService, _ := setupTestController()
	id := "1"
	input := mattermostdto.MattermostNotificationChannelRequest{
		ChannelName: "test",
		WebhookUrl:  "https://webhookurl.com/hooks/id1",
		Description: "desc"}
	updated := models.NotificationChannel{
		Id:          helper.ToPtr(id),
		ChannelName: helper.ToPtr("test"),
		WebhookUrl:  helper.ToPtr("https://webhookurl.com/hooks/id1"),
		Description: helper.ToPtr("desc")}

	notificationService.
		On("UpdateNotificationChannel", mock.Anything, id, mock.Anything).
		Return(updated, nil)

	httpassert.New(t, r).
		Put("/notification-channel/mattermost/1").
		JsonContentObject(input).
		Expect().
		StatusCode(http.StatusOK).
		JsonPath("$.channelName", "test").
		JsonPath("$.webhookUrl", "https://webhookurl.com/hooks/id1").
		JsonPath("$.description", "desc").
		JsonPath("$.id", id)
}

func TestUpdateMattermostChannel_BadRequest(t *testing.T) {
	r, _, _ := setupTestController()
	badBody := []byte(`{"invalid":}`)
	httpassert.New(t, r).
		Put("/notification-channel/mattermost/1").
		JsonContentObject(badBody).
		Expect().
		StatusCode(http.StatusBadRequest)
}

func TestDeleteMattermostChannel_Success(t *testing.T) {
	r, notificationService, _ := setupTestController()
	notificationService.
		On("DeleteNotificationChannel", mock.Anything, "1").
		Return(nil)

	httpassert.New(t, r).
		Delete("/notification-channel/mattermost/1").
		Expect().
		StatusCode(http.StatusNoContent)
}

func TestDeleteMattermostChannel_Error(t *testing.T) {
	r, notificationService, _ := setupTestController()
	notificationService.
		On("DeleteNotificationChannel", mock.Anything, "1").
		Return(errors.New("fail"))

	httpassert.New(t, r).
		Delete("/notification-channel/mattermost/1").
		Expect().
		StatusCode(http.StatusInternalServerError)
}
