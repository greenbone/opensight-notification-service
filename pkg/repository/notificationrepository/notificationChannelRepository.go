package notificationrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/jmoiron/sqlx"
)

type NotificationChannelRepository struct {
	client *sqlx.DB
}

func NewNotificationChannelRepository(db *sqlx.DB) (port.NotificationChannelRepository, error) {
	if db == nil {
		return nil, errors.New("nil db reference")
	}
	client := &NotificationChannelRepository{
		client: db,
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

const updateNotificationChannelQuery = `
        UPDATE notification_service.notification_channel SET
            channel_type = :channel_type,
            channel_name = :channel_name,
            webhook_url = :webhook_url,
            domain = :domain,
            port = :port,
            is_authentication_required = :is_authentication_required,
            is_tls_enforced = :is_tls_enforced,
            username = :username,
            password = :password,
            max_email_attachment_size_mb = :max_email_attachment_size_mb,
            max_email_include_size_mb = :max_email_include_size_mb,
            sender_email_address = :sender_email_address,
            updated = NOW()
        WHERE id = :id
        RETURNING *
    `

func (r *NotificationChannelRepository) CreateNotificationChannel(ctx context.Context, channelIn models.NotificationChannel) (models.NotificationChannel, error) {
	insertRow, err := toNotificationChannelRow(channelIn)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("invalid argument for inserting notification channel: %w", err)
	}

	stmt, err := r.client.PrepareNamedContext(ctx, createNotificationChannelQuery)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not prepare sql statement: %w", err)
	}
	defer stmt.Close()

	var row notificationChannelRow
	err = stmt.QueryRowxContext(ctx, insertRow).StructScan(&row)
	if err != nil {
		return models.NotificationChannel{}, fmt.Errorf("could not insert into database: %w", err)
	}

	channel := row.ToModel()
	return channel, nil
}

func (r *NotificationChannelRepository) ListNotificationChannelsByType(ctx context.Context, channelType string) ([]models.NotificationChannel, error) {
	query := `SELECT * FROM notification_service.notification_channel WHERE channel_type = $1`
	var rows []notificationChannelRow
	err := r.client.SelectContext(ctx, &rows, query, channelType)
	if err != nil {
		return nil, fmt.Errorf("select by type failed: %w", err)
	}
	result := make([]models.NotificationChannel, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToModel())
	}
	return result, nil
}

func (r *NotificationChannelRepository) DeleteNotificationChannel(ctx context.Context, id string) error {
	query := `DELETE FROM notification_service.notification_channel WHERE id = $1`
	_, err := r.client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

// UpdateNotificationChannel updates an existing notification channel and returns the updated model.
func (r *NotificationChannelRepository) UpdateNotificationChannel(ctx context.Context, id string, in models.NotificationChannel) (models.NotificationChannel, error) {
	rowIn, err := toNotificationChannelRow(in)
	if err != nil {
		return in, fmt.Errorf("convert to row failed: %w", err)
	}
	rowIn.Id = id

	var row notificationChannelRow
	stmt, err := r.client.PrepareNamedContext(ctx, updateNotificationChannelQuery)
	if err != nil {
		return in, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, rowIn).StructScan(&row)
	if err != nil {
		return in, fmt.Errorf("update failed: %w", err)
	}

	return row.ToModel(), nil
}
