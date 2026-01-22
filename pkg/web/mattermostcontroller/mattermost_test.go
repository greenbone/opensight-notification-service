package mattermostcontroller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/request"
)

func setupTestController() (*gin.Engine, *mocks.NotificationChannelService, *mocks.MattermostChannelService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	notificationService := &mocks.NotificationChannelService{}
	mattermostService := &mocks.MattermostChannelService{}
	ctrl := &MattermostController{
		service:                  notificationService,
		mattermostChannelService: mattermostService,
	}
	group := r.Group("/notification-channel/mattermost")
	group.POST("", ctrl.CreateMattermostChannel)
	group.GET("", ctrl.ListMattermostChannelsByType)
	group.PUT(":id", ctrl.UpdateMattermostChannel)
	group.DELETE(":id", ctrl.DeleteMattermostChannel)
	return r, notificationService, mattermostService
}

func TestCreateMattermostChannel_Success(t *testing.T) {
	r, _, mattermostService := setupTestController()
	input := request.MattermostNotificationChannelRequest{
		ChannelName: "test-channel",
		WebhookUrl:  "http://webhook",
		Description: "desc",
	}
	mattermostService.On("CreateMattermostChannel", mock.Anything, input).Return(input, nil)

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
	id := strPtr("123")
	desc := strPtr("desc")
	channels := []models.NotificationChannel{{
		Id:          id,
		ChannelName: strPtr("test"),
		WebhookUrl:  strPtr("url"),
		Description: desc,
	}}
	notificationService.On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).Return(channels, nil)

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
	notificationService.On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).Return(nil, errors.New("fail"))

	httpassert.New(t, r).
		Get("/notification-channel/mattermost").
		Expect().
		StatusCode(http.StatusInternalServerError)
}

func TestUpdateMattermostChannel_Success(t *testing.T) {
	r, notificationService, _ := setupTestController()
	id := "1"
	input := request.MattermostNotificationChannelRequest{ChannelName: "test", WebhookUrl: "url", Description: "desc", Id: strPtr(id)}
	updated := models.NotificationChannel{Id: strPtr(id), ChannelName: strPtr("test"), WebhookUrl: strPtr("url"), Description: strPtr("desc")}
	notificationService.On("UpdateNotificationChannel", mock.Anything, id, mock.Anything).Return(updated, nil)

	httpassert.New(t, r).
		Put("/notification-channel/mattermost/1").
		JsonContentObject(input).
		Expect().
		StatusCode(http.StatusOK).
		JsonPath("$.channelName", "test").
		JsonPath("$.webhookUrl", "url").
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
	notificationService.On("DeleteNotificationChannel", mock.Anything, "1").Return(nil)

	httpassert.New(t, r).
		Delete("/notification-channel/mattermost/1").
		Expect().
		StatusCode(http.StatusNoContent)
}

func TestDeleteMattermostChannel_Error(t *testing.T) {
	r, notificationService, _ := setupTestController()
	notificationService.On("DeleteNotificationChannel", mock.Anything, "1").Return(errors.New("fail"))

	httpassert.New(t, r).
		Delete("/notification-channel/mattermost/1").
		Expect().
		StatusCode(http.StatusInternalServerError)
}

func strPtr(s string) *string { return &s }
