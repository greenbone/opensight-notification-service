// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package ruleservice

import (
	"context"
	"testing"

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
			Levels:  []string{"info"},
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
					Levels:  []string{"info"},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMattermost)},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMail)},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "someone@example.com",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMattermost)},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientNotSupported,
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
					Levels:  []string{"info"},
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
					Levels:  []string{"info"},
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
					Levels:  []string{"info"},
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
				Levels:  []string{"info"},
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
				Levels:  []string{"info"},
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
				Levels:  []string{"info"},
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
					Levels:  []string{"info"},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMattermost)},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "",
				},
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMail)},
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
					Levels:  []string{"info"},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"},
					Recipient: "someone@example.com",
				},
				Active: true,
			},
			mockChannelRepoGet: mockChannelGetCall{
				channelID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				channel:   models.NotificationChannel{ChannelType: string(models.ChannelTypeMattermost)},
				err:       nil,
			},
			mockOriginRepoList: mockOriginListCall{
				origins: []entities.Origin{{Class: "test"}},
				err:     nil,
			},
			wantErr: ErrRecipientNotSupported,
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
