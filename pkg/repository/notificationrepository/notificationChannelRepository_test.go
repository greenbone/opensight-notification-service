// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationrepository

import (
	"context"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/helper"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/greenbone/opensight-notification-service/pkg/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (context.Context, port.NotificationChannelRepository) {
	db := pgtesting.NewDB(t)
	repo, err := NewNotificationChannelRepository(db)
	require.NoError(t, err)
	ctx := context.Background()
	return ctx, repo
}

func Test_NotificationChannelRepository_CRUD(t *testing.T) {
	ctx, repo := setupTestRepo(t)

	// Create
	channelIn := models.NotificationChannel{
		ChannelType:              string(models.ChannelTypeMail),
		ChannelName:              helper.ToPtr("Test Channel"),
		WebhookUrl:               nil,
		Domain:                   helper.ToPtr("example.com"),
		Port:                     helper.ToPtr(587),
		IsAuthenticationRequired: helper.ToPtr(true),
		IsTlsEnforced:            helper.ToPtr(true),
		Username:                 helper.ToPtr("user"),
		Password:                 helper.ToPtr("pass"),
		MaxEmailAttachmentSizeMb: helper.ToPtr(10),
		MaxEmailIncludeSizeMb:    helper.ToPtr(5),
		SenderEmailAddress:       helper.ToPtr("sender@example.com"),
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
	updatedIn.ChannelName = helper.ToPtr("Updated Channel")
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

func Test_NotificationChannelRepository_CreateWithMissingRequiredFields(t *testing.T) {
	ctx, repo := setupTestRepo(t)
	invalidChannel := models.NotificationChannel{}
	_, err := repo.CreateNotificationChannel(ctx, invalidChannel)
	assert.Error(t, err, "expected error for missing required fields")
}

func Test_NotificationChannelRepository_UpdateNonExistentChannel(t *testing.T) {
	ctx, repo := setupTestRepo(t)
	nonExistentId := "00000000-0000-0000-0000-000000000000"
	_, err := repo.UpdateNotificationChannel(ctx, nonExistentId, models.NotificationChannel{
		ChannelType: "mail",
	})
	assert.Error(t, err, "expected error for updating non-existent channel")
}

func Test_NotificationChannelRepository_DeleteNonExistentChannel(t *testing.T) {
	ctx, repo := setupTestRepo(t)
	nonExistentId := "00000000-0000-0000-0000-000000000000"
	err := repo.DeleteNotificationChannel(ctx, nonExistentId)
	assert.NoError(t, err, "deleting non-existent channel should not error")
}

func Test_NotificationChannelRepository_ListWithNonExistentChannelType(t *testing.T) {
	ctx, repo := setupTestRepo(t)
	listed, err := repo.ListNotificationChannelsByType(ctx, "nonexistent-type")
	assert.NoError(t, err)
	assert.Len(t, listed, 0, "expected no channels for unknown type")
}
