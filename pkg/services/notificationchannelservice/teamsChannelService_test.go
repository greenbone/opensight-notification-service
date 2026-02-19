package notificationchannelservice

import (
	"context"
	"net/http"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamscontroller/teamsdto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendTeamsTestMessage_PostsToWebhook(t *testing.T) {
	var gotMethod, gotURL string

	svc := &teamsChannelService{
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
	err := svc.SendTeamsTestMessage(webhook)
	require.NoError(t, err)

	assert.Equal(t, "POST", gotMethod)
	assert.Equal(t, webhook, gotURL)
}

func TestSendTeamsTestMessage_ErrorOnTransport(t *testing.T) {
	svc := &teamsChannelService{
		transport: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return nil, ErrTeamsMessageDelivery
			})},
	}

	err := svc.SendTeamsTestMessage("https://example.com:443/workflows/01fa130f2e134641b2cf39d8a710a002")
	require.ErrorContains(t, err, "teams message could not be send")
}

func TestTeamsChannelLimit(t *testing.T) {
	notificationChannelService := mocks.NewNotificationChannelService(t)
	notificationChannelService.EXPECT().ListNotificationChannelsByType(context.Background(), models.ChannelTypeTeams).
		Return([]models.NotificationChannel{
			{},
		}, nil)

	service := NewTeamsChannelService(notificationChannelService, 1, http.DefaultClient)

	_, err := service.CreateTeamsChannel(context.Background(), teamsdto.TeamsNotificationChannelRequest{})
	require.ErrorIs(t, err, ErrTeamsChannelLimitReached)
}
