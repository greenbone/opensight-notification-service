package notificationchannelservice

import (
	"context"
	"errors"
	"time"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/wneessen/go-mail"
)

var (
	ErrCreateMailClient      = errors.New("failed to create mail client")
	ErrMailServerUnreachable = errors.New("mail server is unreachable")
	ErrCreatingMailMessage   = errors.New("failed to create mail message")
	ErrSendingEmail          = errors.New("failed to send mail")
)

type MailService interface {
	// SendMail sends an HTML email
	SendMail(
		ctx context.Context,
		mailServer models.NotificationChannel,
		receiver string,
		subject string,
		body string,
	) error

	// ConnectionCheck checks connection for host, port and TLS settings
	ConnectionCheck(ctx context.Context, mailServer models.NotificationChannel) error
}

type mailService struct {
}

func NewMailService() MailService {
	return &mailService{}
}

func (m *mailService) SendMail(
	ctx context.Context,
	mailServer models.NotificationChannel,
	receiver string,
	subject string,
	body string,
) error {
	client, err := m.createClient(mailServer)
	if err != nil {
		return errors.Join(err, ErrCreateMailClient)
	}

	defer func() {
		_ = client.Close()
	}()

	message := mail.NewMsg()
	if err := message.From(*mailServer.SenderEmailAddress); err != nil {
		return errors.Join(err, ErrCreatingMailMessage)
	}
	if err := message.To(receiver); err != nil {
		return errors.Join(err, ErrCreatingMailMessage)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, body)

	err = client.DialAndSendWithContext(ctx, message)
	if err != nil {
		return errors.Join(err, ErrSendingEmail)
	}

	return nil
}

func (m *mailService) ConnectionCheck(ctx context.Context, mailServer models.NotificationChannel) error {
	client, err := m.createClient(mailServer)
	if err != nil {
		return errors.Join(err, ErrCreateMailClient)
	}

	defer func() {
		_ = client.Close()
	}()

	if err = client.DialWithContext(ctx); err != nil {
		return errors.Join(err, ErrMailServerUnreachable)
	}

	return nil
}

func (m *mailService) createClient(mailServer models.NotificationChannel) (*mail.Client, error) {
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
	return client, err
}
