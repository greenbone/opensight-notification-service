package notificationrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type NotificationChannelRepository struct {
	client         *sqlx.DB
	encryptManager port.EncryptManager
}

func NewNotificationChannelRepository(db *sqlx.DB, encryptService port.EncryptManager) (port.NotificationChannelRepository, error) {
	if db == nil {
		return nil, errors.New("nil db reference")
	}
	client := &NotificationChannelRepository{
		client:         db,
		encryptManager: encryptService,
	}
	return client, nil
}

const createNotificationChannelQuery = `
    INSERT INTO notification_service.notification_channel (
        channel_type, channel_name, webhook_url, domain, port,
        is_authentication_required, is_tls_enforced, username, password,
        max_email_attachment_size_mb, max_email_include_size_mb, sender_email_address 
    ) VALUES (
        :channel_type, :channel_name, :webhook_url, :domain, :port,
        :is_authentication_required, :is_tls_enforced, :username, :password,
        :max_email_attachment_size_mb, :max_email_include_size_mb, :sender_email_address
    )
    RETURNING *
`

func buildUpdateNotificationChannelQuery(in models.NotificationChannel) string {
	query := `
        UPDATE notification_service.notification_channel SET
            channel_type = :channel_type,
            channel_name = :channel_name,
            webhook_url = :webhook_url,
            domain = :domain,
            port = :port,
            is_authentication_required = :is_authentication_required,
            is_tls_enforced = :is_tls_enforced,
            username = :username,`

	if in.Password != nil {
		query += `password = :password,`
	}

	query += `
            max_email_attachment_size_mb = :max_email_attachment_size_mb,
            max_email_include_size_mb = :max_email_include_size_mb,
            sender_email_address = :sender_email_address,
            updated_at = NOW()
        WHERE id = :id
        RETURNING *`
	return query
}

// CreateNotificationChannel is now transactional and supports commit/rollback.
func (r *NotificationChannelRepository) CreateNotificationChannel(
	ctx context.Context,
	channelIn models.NotificationChannel,
) (models.NotificationChannel, error) {
	insertRow := toNotificationChannelRow(channelIn)

	rowWithEncryption, err := r.withEncryptedValues(insertRow)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not encrypt password: %w", err)
	}

	tx, err := r.client.BeginTxx(ctx, nil)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	stmt, err := tx.PrepareNamedContext(ctx, createNotificationChannelQuery)
	if err != nil {
		_ = tx.Rollback()
		return models.NotificationChannel{}, fmt.Errorf("could not prepare sql statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var row notificationChannelRow
	err = stmt.QueryRowxContext(ctx, rowWithEncryption).StructScan(&row)
	if err != nil {
		_ = tx.Rollback()
		return models.NotificationChannel{}, fmt.Errorf("could not insert into database: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return r.withPasswordDecrypted(row).ToModel(), nil
}

func (r *NotificationChannelRepository) GetNotificationChannelByIdAndType(
	ctx context.Context,
	id string,
	channelType models.ChannelType,
) (models.NotificationChannel, error) {
	query := `SELECT * FROM notification_service.notification_channel WHERE id = $1 AND channel_type = $2`

	var row notificationChannelRow
	if err := r.client.GetContext(ctx, &row, query, id, channelType); err != nil {
		return models.NotificationChannel{}, fmt.Errorf("select by id failed: %w", err)
	}

	return row.ToModel(), nil
}

func (r *NotificationChannelRepository) ListNotificationChannelsByType(
	ctx context.Context,
	channelType models.ChannelType,
) ([]models.NotificationChannel, error) {
	query := `SELECT * FROM notification_service.notification_channel WHERE channel_type = $1`

	var rows []notificationChannelRow
	if err := r.client.SelectContext(ctx, &rows, query, string(channelType)); err != nil {
		return nil, fmt.Errorf("select by type failed: %w", err)
	}

	result := make([]models.NotificationChannel, 0, len(rows))
	for _, row := range rows {
		result = append(result, r.withPasswordDecrypted(row).ToModel())
	}

	return result, nil
}

// UpdateNotificationChannel is now transactional.
func (r *NotificationChannelRepository) UpdateNotificationChannel(
	ctx context.Context,
	id string,
	in models.NotificationChannel,
) (models.NotificationChannel, error) {
	rowIn := toNotificationChannelRow(in)
	rowIn.Id = &id

	rowWithEncryption, err := r.withEncryptedValues(rowIn)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not encrypt password: %w", err)
	}

	tx, err := r.client.BeginTxx(ctx, nil)
	if err != nil {
		return in, fmt.Errorf("could not begin transaction: %w", err)
	}

	stmt, err := tx.PrepareNamedContext(ctx, buildUpdateNotificationChannelQuery(in))
	if err != nil {
		return in, fmt.Errorf("prepare failed: %w", err)
	}
	defer func() {
		_ = stmt.Close()
		_ = tx.Rollback()
	}()

	var row notificationChannelRow

	err = stmt.QueryRowxContext(ctx, rowWithEncryption).StructScan(&row)
	if err != nil {
		return in, fmt.Errorf("update failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return in, fmt.Errorf("could not commit transaction: %w", err)
	}

	return r.withPasswordDecrypted(row).ToModel(), nil
}

func (r *NotificationChannelRepository) withEncryptedValues(row notificationChannelRow) (notificationChannelRow, error) {
	if row.Password != nil && strings.TrimSpace(*row.Password) != "" {
		encryptedPasswd, err := r.encryptManager.Encrypt(*row.Password)
		if err != nil {
			return empty, fmt.Errorf("could not encrypt password: %w", err)
		}

		passwd := string(encryptedPasswd)
		row.Password = &passwd
	}

	if row.Username != nil && strings.TrimSpace(*row.Username) != "" {
		encryptedUsername, err := r.encryptManager.Encrypt(*row.Username)
		if err != nil {
			return empty, fmt.Errorf("could not encrypt password: %w", err)
		}

		username := string(encryptedUsername)
		row.Username = &username
	}

	return row, nil
}

func (r *NotificationChannelRepository) withPasswordDecrypted(row notificationChannelRow) notificationChannelRow {
	if row.Password != nil && strings.TrimSpace(*row.Password) != "" {
		dPasswd := *row.Password
		dcPassword, err := r.encryptManager.Decrypt([]byte(dPasswd))
		if err != nil {
			log.Err(err).Msg("could not decrypt password")
		}

		row.Password = &dcPassword
	}

	if row.Username != nil && strings.TrimSpace(*row.Username) != "" {
		username := *row.Username
		dcUsername, err := r.encryptManager.Decrypt([]byte(username))
		if err != nil {
			log.Err(err).Msg("could not decrypt password")
		}

		row.Username = &dcUsername
	}

	return row
}

// DeleteNotificationChannel is now transactional.
func (r *NotificationChannelRepository) DeleteNotificationChannel(ctx context.Context, id string) error {
	tx, err := r.client.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	query := `DELETE FROM notification_service.notification_channel WHERE id = $1`

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err = tx.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
