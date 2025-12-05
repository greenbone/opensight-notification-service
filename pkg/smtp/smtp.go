package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Dialer   func(ctx context.Context, addr string) (net.Conn, error)
}

func (conf Config) Validate() error {
	if conf.Host == "" {
		return errors.New("no SMTP server host specified")
	}
	if conf.Port == "" {
		return errors.New("no SMTP server port specified")
	}
	if conf.Username == "" {
		return errors.New("no SMTP server username specified")
	}
	if conf.Password == "" {
		return errors.New("no SMTP server user password specified")
	}
	return nil
}

func dial(ctx context.Context, conf Config, addr string) (net.Conn, error) {
	if conf.Dialer != nil {
		return conf.Dialer(ctx, addr)
	}
	return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
}

func newClientContext(ctx context.Context, conf Config) (*smtp.Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(conf.Host, conf.Port)
	conn, err := dial(ctx, conf, addr)
	if err != nil {
		return nil, fmt.Errorf("connecting to %q: %w", addr, err)
	}

	// if applicable set global timeout for the entire session
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return nil, fmt.Errorf("setting connection deadlines for %q: %w", addr, err)
		}
	}

	client, err := smtp.NewClient(conn, conf.Host)
	if err != nil {
		return nil, fmt.Errorf("connecting to SMTP server %q: %w", addr, err)
	}
	clientClose := client.Close
	defer func() { _ = clientClose() }()

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: conf.Host}
		if err = client.StartTLS(config); err != nil {
			return nil, fmt.Errorf("creating secure SMTP session to %q: %w", addr, err)
		}
	}

	if ok, _ := client.Extension("AUTH"); !ok {
		return nil, fmt.Errorf("SMTP server %q doesn't support AUTH", addr)
	}
	if err := client.Auth(smtp.PlainAuth("", conf.Username, conf.Password, conf.Host)); err != nil {
		return nil, fmt.Errorf("SMTP server %q authentication: %w", addr, err)
	}

	clientClose = func() error { return nil }
	return client, nil
}

func send(client *smtp.Client, from string, to []string, msg []byte) error {
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}

// SendMail sends mail to MTA server using SMTP protocol.
func SendMail(ctx context.Context, conf Config, from string, to []string, msg []byte) error {
	client, err := newClientContext(ctx, conf)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = send(client, from, to, msg); err != nil {
		return fmt.Errorf("SMTP server %q send mail: %w", net.JoinHostPort(conf.Host, conf.Port), err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("SMTP server %q session termination: %w", net.JoinHostPort(conf.Host, conf.Port), err)
	}
	return nil
}

// ChecConnection checks for SMTP connectivity and credentials validity.
func ChecConnection(ctx context.Context, conf Config) error {
	client, err := newClientContext(ctx, conf)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Noop(); err != nil {
		return fmt.Errorf("SMTP server %q request: %w", net.JoinHostPort(conf.Host, conf.Port), err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("SMTP server %q session termination: %w", net.JoinHostPort(conf.Host, conf.Port), err)
	}
	return nil
}
