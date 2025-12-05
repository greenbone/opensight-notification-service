package smtp_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/greenbone/opensight-notification-service/pkg/smtp"
)

func TestConfigValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		given              smtp.Config
		errorShouldContain string
	}{
		{
			name: "valid",
			given: smtp.Config{
				Host:     "host",
				Port:     "1234",
				Username: "username",
				Password: "password",
			},
		},
		{
			name: "missing-host",
			given: smtp.Config{
				Host:     "",
				Port:     "1234",
				Username: "username",
				Password: "password",
			},
			errorShouldContain: "no SMTP server host specified",
		},
		{
			name: "missing-port",
			given: smtp.Config{
				Host:     "host",
				Port:     "",
				Username: "username",
				Password: "password",
			},
			errorShouldContain: "no SMTP server port specified",
		},
		{
			name: "missing-username",
			given: smtp.Config{
				Host:     "host",
				Port:     "1234",
				Username: "",
				Password: "password",
			},
			errorShouldContain: "no SMTP server username specified",
		},
		{
			name: "missing-password",
			given: smtp.Config{
				Host:     "host",
				Port:     "1234",
				Username: "username",
				Password: "",
			},
			errorShouldContain: "no SMTP server user password specified",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.given.Validate()
			if test.errorShouldContain == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errorShouldContain)
			}
		})
	}
}
