// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package ruleservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/notifications"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/ruleservice/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getValidRule() models.Rule {
	return models.Rule{
		Name: "Valid Rule",
		Trigger: models.Trigger{
			Origins: []models.OriginReference{{Class: "test"}},
			Levels:  []notifications.Level{notifications.LevelInfo},
		},
		Action: models.Action{
			Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
		},
		Active: true,
	}
}

func TestRuleService_Create_RuleLimitReached(t *testing.T) {
	t.Parallel()
	mockRuleRepo := mocks.NewRuleRepository(t)
	mockChannelRepo := mocks.NewNotificationChannelRepository(t)
	mockOriginRepo := mocks.NewOriginRepository(t)

	ruleLimit := 5
	service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, ruleLimit)

	// Mock List to return exactly ruleLimit number of rules
	existingRules := make([]models.Rule, ruleLimit)
	mockRuleRepo.EXPECT().List(mock.Anything).Return(existingRules, nil)

	ctx := context.Background()
	_, err := service.Create(ctx, getValidRule())

	assert.ErrorIs(t, err, ErrRuleLimitReached)
}

type mockChannelGetCall struct {
	channelID string
	channel   models.NotificationChannel
	err       error
}

type mockOriginListCall struct {
	origins []entities.Origin
	err     error
}

type mockChannelListByTypeCall struct {
	channels []models.NotificationChannel
	err      error
}

func strPtr(s string) *string { return &s }

func TestRuleService_Create(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		rule               models.Rule
		mockChannelRepoGet mockChannelGetCall
		mockOriginRepoList mockOriginListCall
		wantErr            error
	}{
		"missing channel": {
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "non-existent-channel"},
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "non-existent-channel",
				err:       errs.ErrItemNotFound,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrChannelNotFound,
		},
		"missing origin": {
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "non-existent-origin"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMattermost},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{
					{Class: "existing-origin-1"},
					{Class: "existing-origin-2"},
				},
				err: nil,
			},
			wantErr: ErrOriginsNotFound,
		},
		"recipient required for mail channel": {
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMail},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientRequired,
		},
		"recipient not supported for non-mail channel": {
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "someone@example.com",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMattermost},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientNotSupported,
		},
		"channel with disallowed type should fail": {
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: "unsupported-channel-type"},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrChannelNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockRuleRepo := mocks.NewRuleRepository(t)
			mockChannelRepo := mocks.NewNotificationChannelRepository(t)
			mockOriginRepo := mocks.NewOriginRepository(t)

			service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, 10)

			mockRuleRepo.EXPECT().List(mock.Anything).Return([]models.Rule{}, nil)
			mockRuleRepo.EXPECT().Create(mock.Anything, tt.rule).Return(models.Rule{}, nil).Maybe()

			mockChannelRepo.EXPECT().GetNotificationChannelById(mock.Anything, tt.mockChannelRepoGet.channelID).
				Return(tt.mockChannelRepoGet.channel, tt.mockChannelRepoGet.err).Once()
			mockOriginRepo.EXPECT().ListOrigins(mock.Anything).Return(tt.mockOriginRepoList.origins, tt.mockOriginRepoList.err).Once()

			ctx := context.Background()
			_, err := service.Create(ctx, tt.rule)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestRuleService_Get_InvalidRuleDeactivated(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		rule            models.Rule
		wantDeactivated bool
		wantErrorField  bool
	}{
		"rule with missing origins is deactivated": {
			rule: models.Rule{
				ID:   "rule-2",
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{}, // empty origins
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
				Active: true,
			},
			wantDeactivated: true,
			wantErrorField:  true,
		},
		"rule with missing channel ID is deactivated": {
			rule: models.Rule{
				ID:   "rule-4",
				Name: "Test Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: ""}, // missing channel ID
				},
				Active: true,
			},
			wantDeactivated: true,
			wantErrorField:  true,
		},
		"valid rule remains active": {
			rule: models.Rule{
				ID:   "rule-5",
				Name: "Valid Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
				Active: true,
			},
			wantDeactivated: false,
			wantErrorField:  false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockRuleRepo := mocks.NewRuleRepository(t)
			mockChannelRepo := mocks.NewNotificationChannelRepository(t)
			mockOriginRepo := mocks.NewOriginRepository(t)

			service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, 10)

			mockRuleRepo.EXPECT().Get(mock.Anything, tt.rule.ID).Return(tt.rule, nil)

			ctx := context.Background()
			got, err := service.Get(ctx, tt.rule.ID)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantDeactivated, !got.Active)
			if tt.wantErrorField {
				assert.NotEmpty(t, got.Errors)
			} else {
				assert.Empty(t, got.Errors)
			}
			got.Errors = tt.rule.Errors // ignore errors field for comparison
			got.Active = tt.rule.Active // ignore active field for comparison
			assert.Equal(t, tt.rule, got)
		})
	}
}

func TestRuleService_List_InvalidRulesDeactivated(t *testing.T) {
	t.Parallel()

	mockRuleRepo := mocks.NewRuleRepository(t)
	mockChannelRepo := mocks.NewNotificationChannelRepository(t)
	mockOriginRepo := mocks.NewOriginRepository(t)

	service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, 10)

	rulesFromRepo := []models.Rule{
		{
			ID:   "rule-1",
			Name: "Valid Rule",
			Trigger: models.Trigger{
				Origins: []models.OriginReference{{Class: "test"}},
				Levels:  []notifications.Level{notifications.LevelInfo},
			},
			Action: models.Action{
				Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
			},
			Active: true,
		},
		{
			ID:   "rule-2",
			Name: "Rule 2",
			Trigger: models.Trigger{
				Origins: []models.OriginReference{}, // invalid - missing origins
				Levels:  []notifications.Level{notifications.LevelInfo},
			},
			Action: models.Action{
				Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
			},
			Active: true,
		},
		{
			ID:   "rule-3",
			Name: "Rule 3",
			Trigger: models.Trigger{
				Origins: []models.OriginReference{{Class: "test"}},
				Levels:  []notifications.Level{notifications.LevelInfo},
			},
			Action: models.Action{
				Channel: models.ChannelReference{}, // invalid - missing channel ID
			},
			Active: true,
		},
	}

	mockRuleRepo.EXPECT().List(mock.Anything).Return(rulesFromRepo, nil)

	ctx := context.Background()
	results, err := service.List(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// First rule should remain active
	assert.True(t, results[0].Active)
	assert.Empty(t, results[0].Errors)

	// Second rule should be deactivated due to missing origin
	assert.False(t, results[1].Active)
	assert.NotEmpty(t, results[1].Errors)

	// Third rule should remain active due to missing channel ID
	assert.False(t, results[2].Active)
	assert.NotEmpty(t, results[2].Errors)

	for i := range results {
		results[i].Errors = rulesFromRepo[i].Errors // ignore errors field for comparison
		results[i].Active = rulesFromRepo[i].Active // ignore active field for comparison
	}
	assert.Equal(t, rulesFromRepo, results)
}

func TestRuleService_Update(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		ruleID             string
		rule               models.Rule
		mockChannelRepoGet mockChannelGetCall
		mockOriginRepoList mockOriginListCall
		wantErr            error
	}{
		"missing channel": {
			ruleID: "rule-1",
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "non-existent-channel"},
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "non-existent-channel",
				err:       errs.ErrItemNotFound,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrChannelNotFound,
		},
		"missing origin": {
			ruleID: "rule-1",
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "non-existent-origin"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMattermost},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{
					{Class: "existing-origin-1"},
					{Class: "existing-origin-2"},
				},
				err: nil,
			},
			wantErr: ErrOriginsNotFound,
		},
		"recipient required for mail channel": {
			ruleID: "rule-1",
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "",
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMail},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientRequired,
		},
		"recipient not supported for non-mail channel": {
			ruleID: "rule-1",
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "someone@example.com",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: models.ChannelTypeMattermost},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientNotSupported,
		},
		"channel with disallowed type should fail": {
			ruleID: "rule-1",
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Origins: []models.OriginReference{{Class: "test"}},
					Levels:  []notifications.Level{notifications.LevelInfo},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: "unsupported-channel-type"},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrChannelNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockRuleRepo := mocks.NewRuleRepository(t)
			mockChannelRepo := mocks.NewNotificationChannelRepository(t)
			mockOriginRepo := mocks.NewOriginRepository(t)

			service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, 10)

			mockRuleRepo.EXPECT().Update(mock.Anything, tt.ruleID, tt.rule).Return(models.Rule{}, nil).Maybe()

			mockChannelRepo.EXPECT().GetNotificationChannelById(mock.Anything, tt.mockChannelRepoGet.channelID).
				Return(tt.mockChannelRepoGet.channel, tt.mockChannelRepoGet.err).Once()
			mockOriginRepo.EXPECT().ListOrigins(mock.Anything).Return(tt.mockOriginRepoList.origins, tt.mockOriginRepoList.err).Once()

			ctx := context.Background()
			_, err := service.Update(ctx, tt.ruleID, tt.rule)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestRuleService_GetAllRuleOptionsFiltered(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		mockOriginRepoList    mockOriginListCall
		mockChannelListByType map[models.ChannelType]mockChannelListByTypeCall
		wantErr               bool
		wantOriginCount       int
		wantChannelCount      int
		wantLevels            []notifications.Level
	}{
		"returns origins, levels and channels of different types": {
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{
					{Class: "origin-1"},
					{Class: "origin-2"},
				},
				err: nil,
			},
			mockChannelListByType: map[models.ChannelType]mockChannelListByTypeCall{
				models.ChannelTypeMail: {
					channels: []models.NotificationChannel{
						{ChannelType: string(models.ChannelTypeMail), ChannelName: strPtr("Mail Channel 1")},
					},
				},
				models.ChannelTypeMattermost: {
					channels: []models.NotificationChannel{
						{ChannelType: string(models.ChannelTypeMattermost), ChannelName: strPtr("Mattermost Channel 1")},
						{ChannelType: string(models.ChannelTypeMattermost), ChannelName: strPtr("Mattermost Channel 2")},
					},
				},
				models.ChannelTypeTeams: {
					channels: []models.NotificationChannel{
						{ChannelType: string(models.ChannelTypeTeams), ChannelName: strPtr("Teams Channel 1")},
					},
				},
			},
			wantErr:          false,
			wantOriginCount:  2,
			wantChannelCount: 4,
			wantLevels:       notifications.AllowedLevels,
		},
		"returns empty origins and no channels": {
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{},
				err:     nil,
			},
			mockChannelListByType: map[models.ChannelType]mockChannelListByTypeCall{
				models.ChannelTypeMail:       {channels: []models.NotificationChannel{}},
				models.ChannelTypeMattermost: {channels: []models.NotificationChannel{}},
				models.ChannelTypeTeams:      {channels: []models.NotificationChannel{}},
			},
			wantErr:          false,
			wantOriginCount:  0,
			wantChannelCount: 0,
			wantLevels:       notifications.AllowedLevels,
		},
		"only mattermost channels configured": {
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "origin-1"}},
				err:     nil,
			},
			mockChannelListByType: map[models.ChannelType]mockChannelListByTypeCall{
				models.ChannelTypeMail: {channels: []models.NotificationChannel{}},
				models.ChannelTypeMattermost: {
					channels: []models.NotificationChannel{
						{ChannelType: string(models.ChannelTypeMattermost), ChannelName: strPtr("Mattermost Only")},
					},
				},
				models.ChannelTypeTeams: {channels: []models.NotificationChannel{}},
			},
			wantErr:          false,
			wantOriginCount:  1,
			wantChannelCount: 1,
			wantLevels:       notifications.AllowedLevels,
		},
		"origin repo error": {
			mockOriginRepoList: mockOriginListCall{
				origins: nil,
				err:     fmt.Errorf("database error"),
			},
			mockChannelListByType: nil,
			wantErr:               true,
		},
		"channel repo error for mail type": {
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "origin-1"}},
				err:     nil,
			},
			mockChannelListByType: map[models.ChannelType]mockChannelListByTypeCall{
				models.ChannelTypeMail: {err: fmt.Errorf("mail channel db error")},
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockRuleRepo := mocks.NewRuleRepository(t)
			mockChannelRepo := mocks.NewNotificationChannelRepository(t)
			mockOriginRepo := mocks.NewOriginRepository(t)

			service := NewRuleService(mockRuleRepo, mockChannelRepo, mockOriginRepo, 10)

			mockOriginRepo.EXPECT().ListOrigins(mock.Anything).Return(tt.mockOriginRepoList.origins, tt.mockOriginRepoList.err).Once()

			// Only set up channel mocks if origin listing succeeds (channels are fetched after origins)
			if tt.mockOriginRepoList.err == nil {
				for _, channelType := range []models.ChannelType{models.ChannelTypeMail, models.ChannelTypeMattermost, models.ChannelTypeTeams} {
					if call, ok := tt.mockChannelListByType[channelType]; ok {
						mockChannelRepo.EXPECT().ListNotificationChannelsByType(mock.Anything, channelType).
							Return(call.channels, call.err).Once()
						if call.err != nil {
							break // service stops iterating on first error
						}
					}
				}
			}

			ctx := context.Background()
			result, err := service.GetAllRuleOptionsFiltered(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Origins, tt.wantOriginCount)
				assert.Len(t, result.Channels, tt.wantChannelCount)
				assert.Equal(t, tt.wantLevels, result.Levels)
				assert.Equal(t, tt.mockOriginRepoList.origins, result.Origins)
			}
		})
	}
}
