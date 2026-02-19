package notificationchannelservice

import (
	"context"
	"net/http"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller/mattermostdto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMattermostTestMessage_PostsToWebhook(t *testing.T) {
	var gotMethod, gotURL string

	svc := &mattermostChannelService{
		transport: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				gotMethod = r.Method
				gotURL = r.URL.String()
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
					Header:     make(http.Header),
				}, nil
			})},
	}

	webhook := "https://example.com:443/workflows/01fa130f2e134641b2cf39d8a710a002"
	err := svc.SendMattermostTestMessage(webhook)
	require.NoError(t, err)

	assert.Equal(t, "POST", gotMethod)
	assert.Equal(t, webhook, gotURL)
}

func TestSendMattermostTestMessage_ErrorOnTransport(t *testing.T) {
	svc := &mattermostChannelService{
		transport: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return nil, ErrMattermostMassageDelivery
			})},
	}

	err := svc.SendMattermostTestMessage("https://example.com:443/workflows/01fa130f2e134641b2cf39d8a710a002")
	require.ErrorContains(t, err, "mattermost message could not be send")
}

func TestMattermostChannelLimit(t *testing.T) {
	notificationChannelService := mocks.NewNotificationChannelService(t)
	notificationChannelService.EXPECT().ListNotificationChannelsByType(context.Background(), models.ChannelTypeMattermost).
		Return([]models.NotificationChannel{
			{},
		}, nil)

	service := NewMattermostChannelService(notificationChannelService, 1, http.DefaultClient)

	_, err := service.CreateMattermostChannel(context.Background(), mattermostdto.MattermostNotificationChannelRequest{})
	require.ErrorIs(t, err, ErrMattermostChannelLimitReached)
}
