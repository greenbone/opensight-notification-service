package notificationchannelservice

import (
	"errors"
	"net/http"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendTeamsTestMessage_PostsToWebhook(t *testing.T) {
	var gotMethod, gotURL string

	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		gotMethod = r.Method
		gotURL = r.URL.String()
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}, nil
	})

	client := http.Client{Transport: rt}
	svc := &teamsChannelService{
		transport: client,
	}

	webhook := "https://default41af85004428462ca2334f3ae673f7.bb.environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/01fa130f2e134641b2cf39d8a710a002/triggers/manual/paths/invoke"
	err := svc.SendTeamsTestMessage(webhook)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotURL != webhook {
		t.Errorf("expected URL %s, got %s", webhook, gotURL)
	}
}

func TestSendTeamsTestMessage_ErrorOnTransport(t *testing.T) {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("transport error")
	})
	client := http.Client{Transport: rt}
	svc := &teamsChannelService{
		transport: client,
	}
	err := svc.SendTeamsTestMessage("https://default41af85004428462ca2334f3ae673f7.bb.environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/01fa130f2e134641b2cf39d8a710a002/triggers/manual/paths/invoke")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
