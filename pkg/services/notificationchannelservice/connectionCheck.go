package notificationchannelservice

import (
	"context"
	"errors"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/wneessen/go-mail"
)

var (
	ErrCreateMailFailed      = errors.New("failed to create mail client")
	ErrMailServerUnreachable = errors.New("mail server is unreachable")
)

func ConnectionCheckMail(ctx context.Context, mailServer models.NotificationChannel) error {
	options := []mail.Option{
		mail.WithPort(*mailServer.Port),
		mail.WithTimeout(5 * time.Second),
	}

	if mailServer.IsTlsEnforced != nil && *mailServer.IsTlsEnforced {
		options = append(options, mail.WithSSL())
	}

	if mailServer.IsAuthenticationRequired != nil && *mailServer.IsAuthenticationRequired {
		options = append(options, mail.WithUsername(*mailServer.Username), mail.WithPassword(*mailServer.Password))
	}

	client, err := mail.NewClient(
		*mailServer.Domain,
		options...,
	)
	defer func() {
		_ = client.Close()
	}()

	if err != nil {
		return errors.Join(err, ErrCreateMailFailed)
	}

	if err = client.DialWithContext(ctx); err != nil {
		return errors.Join(err, ErrMailServerUnreachable)
	}

	// TODO: 21.01.2026 stolksdorf - username and password are not validated

	return nil
}
