package notificationchannelservice

import (
	"net/http"
	"testing"

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
		transport: http.Client{
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
		transport: http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return nil, assert.AnError
			})},
	}

	err := svc.SendTeamsTestMessage("https://example.com:443/workflows/01fa130f2e134641b2cf39d8a710a002")
	require.ErrorContains(t, err, "can not send teams test message")
}
