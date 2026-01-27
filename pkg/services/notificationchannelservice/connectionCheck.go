package notificationchannelservice

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/wneessen/go-mail"
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
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err = client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("failed to reach mail server: %w", err)
	}

	return nil
}
