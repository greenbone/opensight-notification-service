// SPDX-FileCopyrightText: 2026 Danny Ward <https://github.com/dannyward630>
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationchannelservice

import (
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestMailServiceCreateClientRejectsMissingAuthenticationCredentials(t *testing.T) {
	empty := ""
	tests := []struct {
		name     string
		username *string
		password *string
	}{
		{name: "missing both credentials"},
		{name: "missing username", password: new("password")},
		{name: "missing password", username: new("username")},
		{name: "empty username", username: &empty, password: new("password")},
		{name: "empty password", username: new("username"), password: &empty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := (&mailService{}).createClient(models.NotificationChannel{
				Domain:                   new("127.0.0.1"),
				Port:                     new(25),
				IsAuthenticationRequired: new(true),
				Username:                 tt.username,
				Password:                 tt.password,
			})
			require.ErrorIs(t, err, errMailAuthenticationCredentialsRequired)
		})
	}
}
