package smtp_test

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/greenbone/opensight-notification-service/pkg/smtp"
)

func expectRead(expected string) func(*testing.T, net.Conn) {
	return func(t *testing.T, conn net.Conn) {
		assert.NoError(t, conn.SetReadDeadline(time.Now().Add(2*time.Second)))
		read := ""
		nextChar := make([]byte, 1)
		for read != expected {
			_, err := conn.Read(nextChar)
			if err != nil {
				t.Errorf("unexpected error: %v\ndata: %q", err, read)
				return
			}
			read += string(nextChar)
		}
	}
}

func expectReadPrefix(prefix, readUntil string) func(*testing.T, net.Conn) {
	return func(t *testing.T, conn net.Conn) {
		assert.NoError(t, conn.SetReadDeadline(time.Now().Add(2*time.Second)))
		read := ""
		nextChar := make([]byte, 1)
		for !strings.HasSuffix(read, readUntil) {
			_, err := conn.Read(nextChar)
			if err != nil {
				t.Errorf("unexpected error: %v\ndata: %q", err, read)
				return
			}
			read += string(nextChar)
		}
		if !strings.HasPrefix(read, prefix) {
			t.Errorf("missing prefix %q\ndata: %q", prefix, read)
			return
		}
	}
}

func expectWrite(data string) func(*testing.T, net.Conn) {
	return func(t *testing.T, conn net.Conn) {
		assert.NoError(t, conn.SetWriteDeadline(time.Now().Add(2*time.Second)))
		_, err := conn.Write([]byte(data))
		assert.NoError(t, err)
	}
}

func testTcpListener(t *testing.T, handlers ...func(*testing.T, net.Conn)) net.Addr {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, lis.Close())
	})

	handleConn := func(conn net.Conn) {
		for _, h := range handlers {
			h(t, conn)
		}
		_ = conn.Close()
	}

	acceptConn := func() {
		conn, err := lis.Accept()
		assert.NoError(t, err)
		if conn != nil {
			t.Cleanup(func() {
				_ = conn.Close()
			})
			handleConn(conn)
		}
	}

	go acceptConn()
	return lis.Addr()
}

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

func TestSendMail(t *testing.T) {
	tests := []struct {
		name               string
		startServer        func(t *testing.T) net.Addr
		errorShouldContain string
	}{
		{
			name: "successful-send",
			startServer: func(t *testing.T) net.Addr {
				return testTcpListener(t,
					expectWrite("220 smtp.example.com ESMTP Postfix\r\n"),
					expectReadPrefix("EHLO", "\r\n"),
					expectWrite("250-smtp.example.com Hello localhost [127.0.0.1], pleased to meet you\r\n"),
					expectWrite("250-AUTH PLAIN LOGIN\r\n"),
					expectWrite("250-SIZE 10240000\r\n"),
					expectWrite("250-8BITMIME\r\n"),
					expectWrite("250 ENHANCEDSTATUSCODES\r\n"),
					expectReadPrefix("AUTH PLAIN", "\r\n"),
					expectWrite("235 2.7.0 Authentication successful\r\n"),
					expectReadPrefix("MAIL FROM", "\r\n"),
					expectWrite("250 2.1.0 Ok\r\n"),
					expectReadPrefix("RCPT TO", "\r\n"),
					expectWrite("250 2.1.5 Ok\r\n"),
					expectReadPrefix("DATA", "\r\n"),
					expectWrite("354 Start mail input; end with <CRLF>.<CRLF>\r\n"),
					expectReadPrefix("", "\r\n.\r\n"),
					expectWrite("250 2.0.0 Ok: queued as 12345\r\n"),
					expectReadPrefix("QUIT", "\r\n"),
					expectWrite("221 2.0.0 Bye\r\n"),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr := test.startServer(t)
			host, port, err := net.SplitHostPort(addr.String())
			require.NoError(t, err)

			conf := smtp.Config{
				Host:     host,
				Port:     port,
				Username: "username",
				Password: "password",
			}

			smtpErr := smtp.SendMail(t.Context(), conf, "alice@example.com", []string{"bob@example.com"}, []byte("Message..."))
			if test.errorShouldContain == "" {
				require.NoError(t, smtpErr)
			} else {
				require.Error(t, smtpErr)
				require.Contains(t, smtpErr.Error(), test.errorShouldContain)
			}
		})
	}
}
