// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"context"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NotificationChannelRepository_CRUD(t *testing.T) {
	db := pgtesting.NewDB(t)
	repo, err := NewNotificationChannelRepository(db)
	require.NoError(t, err)
	ctx := context.Background()

	// Create
	channelIn := models.NotificationChannel{
		ChannelType:              "mail",
		ChannelName:              ptrString("Test Channel"),
		WebhookUrl:               nil,
		Domain:                   ptrString("example.com"),
		Port:                     ptrInt(587),
		IsAuthenticationRequired: ptrBool(true),
		IsTlsEnforced:            ptrBool(true),
		Username:                 ptrString("user"),
		Password:                 ptrString("pass"),
		MaxEmailAttachmentSizeMb: ptrInt(10),
		MaxEmailIncludeSizeMb:    ptrInt(5),
		SenderEmailAddress:       ptrString("sender@example.com"),
	}

	created, err := repo.CreateNotificationChannel(ctx, channelIn)
	require.NoError(t, err)
	assert.NotEmpty(t, created.Id)
	assert.Equal(t, channelIn.ChannelType, created.ChannelType)
	assert.Equal(t, *channelIn.ChannelName, *created.ChannelName)

	// List by type
	listed, err := repo.ListNotificationChannelsByType(ctx, "mail")
	require.NoError(t, err)
	assert.Len(t, listed, 1)
	assert.Equal(t, created.Id, listed[0].Id)

	// Update
	updatedIn := created
	updatedIn.ChannelName = ptrString("Updated Channel")
	updated, err := repo.UpdateNotificationChannel(ctx, *created.Id, updatedIn)
	require.NoError(t, err)
	assert.Equal(t, "Updated Channel", *updated.ChannelName)

	// Delete
	err = repo.DeleteNotificationChannel(ctx, *created.Id)
	require.NoError(t, err)

	// List after delete
	listedAfterDelete, err := repo.ListNotificationChannelsByType(ctx, "mail")
	require.NoError(t, err)
	assert.Len(t, listedAfterDelete, 0)
}

func Test_NotificationChannelRepository_NegativeAndEdgeCases(t *testing.T) {
	db := pgtesting.NewDB(t)
	repo, err := NewNotificationChannelRepository(db)
	require.NoError(t, err)
	ctx := context.Background()

	// Create with missing required fields
	invalidChannel := models.NotificationChannel{}
	_, err = repo.CreateNotificationChannel(ctx, invalidChannel)
	assert.Error(t, err, "expected error for missing required fields")

	// Update non-existent channel
	nonExistentId := "00000000-0000-0000-0000-000000000000"
	_, err = repo.UpdateNotificationChannel(ctx, nonExistentId, models.NotificationChannel{
		ChannelType: "mail",
	})
	assert.Error(t, err, "expected error for updating non-existent channel")

	// Delete non-existent channel
	err = repo.DeleteNotificationChannel(ctx, nonExistentId)
	assert.NoError(t, err, "deleting non-existent channel should not error")

	// List with non-existent channelType
	listed, err := repo.ListNotificationChannelsByType(ctx, "nonexistent-type")
	assert.NoError(t, err)
	assert.Len(t, listed, 0, "expected no channels for unknown type")
}

// Helper functions for pointer values
func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrBool(b bool) *bool       { return &b }
