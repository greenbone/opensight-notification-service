package mattermostcontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/request"
)

func setupTestController() (*gin.Engine, *mocks.NotificationChannelService, *mocks.MattermostChannelService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockNotificationService := &mocks.NotificationChannelService{}
	mockMattermostService := &mocks.MattermostChannelService{}
	ctrl := &MattermostController{
		Service:                  mockNotificationService,
		MattermostChannelService: mockMattermostService,
	}
	group := r.Group("/notification-channel/mattermost")
	group.POST("", ctrl.CreateMattermostChannel)
	group.GET("", ctrl.ListMattermostChannelsByType)
	group.PUT(":id", ctrl.UpdateMattermostChannel)
	group.DELETE(":id", ctrl.DeleteMattermostChannel)
	return r, mockNotificationService, mockMattermostService
}

func TestCreateMattermostChannel_Success(t *testing.T) {
	r, _, mockMattermostService := setupTestController()
	input := request.MattermostNotificationChannelRequest{
		ChannelName: "test-channel",
		WebhookUrl:  "http://webhook",
		Description: "desc",
	}
	mockMattermostService.On("CreateMattermostChannel", mock.Anything, input).Return(input, nil)
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/notification-channel/mattermost", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp request.MattermostNotificationChannelRequest
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, input.ChannelName, resp.ChannelName)
}

func TestCreateMattermostChannel_BadRequest(t *testing.T) {
	r, _, _ := setupTestController()
	badBody := []byte(`{"invalid":}`)
	req, _ := http.NewRequest("POST", "/notification-channel/mattermost", bytes.NewBuffer(badBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListMattermostChannelsByType_Success(t *testing.T) {
	r, mockNotificationService, _ := setupTestController()
	id := strPtr("123")
	desc := strPtr("desc")
	channels := []models.NotificationChannel{{
		Id:          id,
		ChannelName: strPtr("test"),
		WebhookUrl:  strPtr("url"),
		Description: desc,
	}}
	mockNotificationService.On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).Return(channels, nil)

	req, _ := http.NewRequest("GET", "/notification-channel/mattermost", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []request.MattermostNotificationChannelRequest
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "test", resp[0].ChannelName)
	assert.Equal(t, "url", resp[0].WebhookUrl)
	assert.Equal(t, "desc", resp[0].Description)
	assert.Equal(t, "123", *resp[0].Id)
}

func TestListMattermostChannelsByType_Error(t *testing.T) {
	r, mockNotificationService, _ := setupTestController()
	mockNotificationService.On("ListNotificationChannelsByType", mock.Anything, models.ChannelTypeMattermost).Return(nil, errors.New("fail"))
	req, _ := http.NewRequest("GET", "/notification-channel/mattermost", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateMattermostChannel_Success(t *testing.T) {
	r, mockNotificationService, _ := setupTestController()
	id := "1"
	input := request.MattermostNotificationChannelRequest{ChannelName: "test", WebhookUrl: "url", Description: "desc", Id: strPtr(id)}
	updated := models.NotificationChannel{Id: strPtr(id), ChannelName: strPtr("test"), WebhookUrl: strPtr("url"), Description: strPtr("desc")}
	mockNotificationService.On("UpdateNotificationChannel", mock.Anything, id, mock.Anything).Return(updated, nil)
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("PUT", "/notification-channel/mattermost/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp request.MattermostNotificationChannelRequest
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test", resp.ChannelName)
	assert.Equal(t, "url", resp.WebhookUrl)
	assert.Equal(t, "desc", resp.Description)
	assert.Equal(t, id, *resp.Id)
}

func TestUpdateMattermostChannel_BadRequest(t *testing.T) {
	r, _, _ := setupTestController()
	badBody := []byte(`{"invalid":}`)
	req, _ := http.NewRequest("PUT", "/notification-channel/mattermost/1", bytes.NewBuffer(badBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteMattermostChannel_Success(t *testing.T) {
	r, mockNotificationService, _ := setupTestController()
	mockNotificationService.On("DeleteNotificationChannel", mock.Anything, "1").Return(nil)
	req, _ := http.NewRequest("DELETE", "/notification-channel/mattermost/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteMattermostChannel_Error(t *testing.T) {
	r, mockNotificationService, _ := setupTestController()
	mockNotificationService.On("DeleteNotificationChannel", mock.Anything, "1").Return(errors.New("fail"))
	req, _ := http.NewRequest("DELETE", "/notification-channel/mattermost/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func strPtr(s string) *string { return &s }
