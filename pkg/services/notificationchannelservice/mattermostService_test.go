// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationchannelservice

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMattermostMessage(t *testing.T) {
	var gotMethod, gotURL string

	svc := &MattermostService{
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
	err := svc.SendMessage(webhook, "test message")
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, webhook, gotURL)
}
