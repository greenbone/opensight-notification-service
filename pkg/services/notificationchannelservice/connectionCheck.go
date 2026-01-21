package notificationchannelservice

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/wneessen/go-mail"
)

func ConnectionCheckMail(ctx context.Context, channel models.NotificationChannel) error {
	options := []mail.Option{
		mail.WithPort(*channel.Port),
		mail.WithTimeout(5 * time.Second),
	}

	if *channel.IsTlsEnforced {
		options = append(options, mail.WithSSL())
	}

	if *channel.IsAuthenticationRequired {
		options = append(options, mail.WithUsername(*channel.Username), mail.WithPassword(*channel.Password))
	}

	client, err := mail.NewClient(
		*channel.Domain,
		options...,
	)
	defer client.Close()

	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err = client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("failed to reach mail server: %w", err)
	}

	// TODO: 21.01.2026 stolksdorf - username and password are not validated

	return nil
}
