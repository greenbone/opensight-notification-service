// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package rulerepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/originrepository"
	"github.com/greenbone/opensight-notification-service/pkg/security"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions

func createTestChannel(t *testing.T, db *sqlx.DB, name, channelType string) (channelID string) {
	// create channel repo
	encryptMgr := security.NewEncryptManager()
	encryptMgr.UpdateKeys(config.DatabaseEncryptionKey{
		Password:     "password",
		PasswordSalt: "password-salt-should-no-be-short-fyi",
	})
	repo, err := notificationrepository.NewNotificationChannelRepository(db, encryptMgr)
	require.NoError(t, err)
	ctx := context.Background()

	// create channel
	channel, err := repo.CreateNotificationChannel(ctx, models.NotificationChannel{
		ChannelName: &name,
		ChannelType: channelType,
	})
	require.NoError(t, err)
	require.NotNil(t, channel.Id)

	return *channel.Id
}

func createTestOrigin(t *testing.T, db *sqlx.DB, name, class, serviceID string) {
	originRepo, err := originrepository.NewOriginRepository(db)
	require.NoError(t, err)

	origins := []entities.Origin{
		{
			Name:      name,
			Class:     class,
			ServiceID: serviceID,
		},
	}
	err = originRepo.UpsertOrigins(context.Background(), serviceID, origins)
	require.NoError(t, err)
}

func Test_GetRule_NotFound(t *testing.T) {
	t.Parallel()
	db := pgtesting.NewDB(t)
	repo, err := NewRuleRepository(db)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = repo.Get(ctx, uuid.NewString())
	assert.ErrorIs(t, err, errs.ErrItemNotFound)
}

func Test_CreateRule_GetRule(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setupData func(t *testing.T, db *sqlx.DB) (channelID string)
		rule      models.Rule
		wantRule  models.Rule
		wantErr   error // if we get an error in this test we always expect it to be a validation error
	}{
		"create rule with single origin and level": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := createTestChannel(t, db, "test-channel", "mattermost")
				createTestOrigin(t, db, "Origin1", "class1", "service1")
				return channelID
			},
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{Class: "class1", Name: "read-only,ignored", ServiceID: "read-only,ignored"},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test", Name: "read-only,ignored", Type: "read-only,ignored"},
				},
				Active: true,
			},
			wantRule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin1",
							Class:     "class1",
							ServiceID: "service1",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{
						ID:   "set below in test",
						Name: "test-channel",
						Type: "mattermost",
					},
				},
				Active: true,
			},
		},
		"create rule with recipient and multiple origins and levels": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := createTestChannel(t, db, "test-channel", "mail")
				createTestOrigin(t, db, "Vulnerability", "vuln", "service1")
				createTestOrigin(t, db, "Compliance", "compliance", "service2")
				return channelID
			},
			rule: models.Rule{
				Name: "Security Alerts",
				Trigger: models.Trigger{
					Levels: []string{"high", "critical"},
					Origins: []models.OriginReference{
						{Class: "vuln", ServiceID: "read-only,ignored", Name: "read-only,ignored"},
						{Class: "compliance", ServiceID: "read-only,ignored", Name: "read-only,ignored"},
					},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "set below in test", Name: "read-only,ignored", Type: "read-only,ignored"},
					Recipient: "security@example.com",
				},
				Active: true,
			},
			wantRule: models.Rule{
				Name: "Security Alerts",
				Trigger: models.Trigger{
					Levels: []string{"high", "critical"},
					Origins: []models.OriginReference{
						{
							Name:      "Vulnerability",
							Class:     "vuln",
							ServiceID: "service1",
						},
						{
							Name:      "Compliance",
							Class:     "compliance",
							ServiceID: "service2",
						},
					},
				},
				Action: models.Action{
					Recipient: "security@example.com",
					Channel: models.ChannelReference{
						ID:   "set below in test",
						Name: "test-channel",
						Type: "mail",
					},
				},
				Active: true,
			},
		},
		"create deactivated rule": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := createTestChannel(t, db, "test-channel", "mattermost")
				createTestOrigin(t, db, "Origin1", "class1", "service1")
				return channelID
			},
			rule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{Class: "class1", Name: "read-only,ignored", ServiceID: "read-only,ignored"},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test", Name: "read-only,ignored", Type: "read-only,ignored"},
				},
				Active: false,
			},
			wantRule: models.Rule{
				Name: "Test Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin1",
							Class:     "class1",
							ServiceID: "service1",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{
						ID:   "",
						Name: "test-channel",
						Type: "mattermost",
					},
				},
				Active: false,
			},
		},
		"create rule with non-existent origin should fail": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := createTestChannel(t, db, "test-channel", "mail")
				return channelID
			},
			rule: models.Rule{
				Name: "Invalid Rule",
				Trigger: models.Trigger{
					Levels:  []string{"medium"},
					Origins: []models.OriginReference{{Class: "non-existent"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test"},
				},
				Active: true,
			},
			wantErr: ErrOriginsNotFound,
		},
		"create rule with non-existent channel works, but returns an empty channel ID": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := uuid.NewString() // non-existent channel ID
				createTestOrigin(t, db, "name1", "class1", "service1")
				return channelID
			},
			rule: models.Rule{
				Name: "Invalid Rule",
				Trigger: models.Trigger{
					Levels:  []string{"medium"},
					Origins: []models.OriginReference{{Class: "class1"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test"},
				},
				Active: true,
			},
			wantRule: models.Rule{
				Name: "Invalid Rule",
				Trigger: models.Trigger{
					Levels: []string{"medium"},
					Origins: []models.OriginReference{
						{
							Name:      "name1",
							Class:     "class1",
							ServiceID: "service1",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{}, // channel not found, so empty channel reference expected
				},
				Active: true,
			},
		},
		"create rule with duplicate name should fail": {
			setupData: func(t *testing.T, db *sqlx.DB) string {
				channelID := createTestChannel(t, db, "test-channel", "mail")
				createTestOrigin(t, db, "name1", "class1", "ns1")
				repo, err := NewRuleRepository(db)
				require.NoError(t, err)
				existingRule := models.Rule{
					Name: "Existing Rule",
					Trigger: models.Trigger{
						Levels:  []string{"medium"},
						Origins: []models.OriginReference{{Class: "class1"}},
					},
					Action: models.Action{
						Channel: models.ChannelReference{ID: channelID},
					},
				}
				_, err = repo.Create(context.Background(), existingRule)
				require.NoError(t, err)

				return channelID
			},
			rule: models.Rule{
				Name: "Existing Rule",
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "class1"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test"},
				},
			},
			wantErr: ErrDuplicateRuleName,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := pgtesting.NewDB(t)

			repo, err := NewRuleRepository(db)
			require.NoError(t, err)

			ctx := context.Background()

			// Setup test data
			channelID := tt.setupData(t, db)
			// set channel in rule
			tt.rule.Action.Channel.ID = channelID

			// Create rule
			createdRule, err := repo.Create(ctx, tt.rule)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			assert.NotEmpty(t, createdRule.ID) // Verify created rule has ID
			tt.wantRule.ID = createdRule.ID    // set id for comparison (not known beforehand)
			if tt.wantRule.Action.Channel.ID != "" {
				tt.wantRule.Action.Channel.ID = channelID // set id for comparison (not known beforehand)
			}

			assert.Equal(t, tt.wantRule, createdRule)

			// Retrieve rule and verify
			gotRule, err := repo.Get(ctx, createdRule.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRule, gotRule)
		})
	}
}

func Test_UpdateRule_NotFound(t *testing.T) {
	t.Parallel()
	db := pgtesting.NewDB(t)
	repo, err := NewRuleRepository(db)
	require.NoError(t, err)
	createTestOrigin(t, db, "test-origin", "class1", "service1")
	channelID := createTestChannel(t, db, "test-channel", "mattermost")

	rule := models.Rule{
		Name: "Non-existent Rule",
		Trigger: models.Trigger{
			Levels:  []string{"low"},
			Origins: []models.OriginReference{{Class: "class1"}},
		},
		Action: models.Action{
			Channel: models.ChannelReference{ID: channelID},
		},
		Active: false,
	}

	ctx := context.Background()
	_, err = repo.Update(ctx, uuid.NewString(), rule)
	assert.ErrorIs(t, err, errs.ErrItemNotFound)
}

func Test_UpdateRule(t *testing.T) {
	t.Parallel()

	origin1 := models.OriginReference{Name: "Origin1", Class: "class1", ServiceID: "service1"}
	origin2 := models.OriginReference{Name: "Origin2", Class: "class2", ServiceID: "service2"}

	channel1 := models.ChannelReference{ID: "channel1", Name: "Channel 1", Type: "mattermost"}
	channel2 := models.ChannelReference{ID: "channel2", Name: "Channel 2", Type: "mail"}

	existingUntouchedRule := models.Rule{
		Name: "Existing untouched Rule",
		Trigger: models.Trigger{
			Levels:  []string{"low"},
			Origins: []models.OriginReference{{Class: "class1"}},
		},
		Action: models.Action{
			Channel: models.ChannelReference{ID: "set below in test", Name: channel1.Name},
		},
		Active: false,
	}
	existingRule := models.Rule{
		Name: "Existing Rule",
		Trigger: models.Trigger{
			Levels:  []string{"medium"},
			Origins: []models.OriginReference{{Class: "class1", Name: "read-only,ignored", ServiceID: "read-only,ignored"}},
		},
		Action: models.Action{
			Channel: models.ChannelReference{ID: "set below in test", Name: channel1.Name},
		},
		Active: false,
	}

	// default setup, create two channels and oigins along with existing rules
	setupData := func(t *testing.T, db *sqlx.DB, repo *RuleRepository) (channelID string, ruleID string) {
		// Setup test data
		channelID1 := createTestChannel(t, db, channel1.Name, channel1.Type)
		channelID2 := createTestChannel(t, db, channel2.Name, channel2.Type)

		createTestOrigin(t, db, origin1.Name, origin1.Class, origin1.ServiceID)
		createTestOrigin(t, db, origin2.Name, origin2.Class, origin2.ServiceID)

		ctx := context.Background()

		// Create existing rules
		existingRule.Action.Channel.ID = channelID1
		createdRule, err := repo.Create(ctx, existingRule)
		require.NoError(t, err)
		existingUntouchedRule.Action.Channel.ID = channelID1
		_, err = repo.Create(ctx, existingUntouchedRule)
		require.NoError(t, err)

		return channelID2, createdRule.ID
	}

	// Note: as the channel ID is not known beforehand, it is set to the one returned by `setupData`
	tests := map[string]struct {
		// setupData returns the channelID which should be used for the updated `Action.Channel.ID` value of the rule
		setupData func(t *testing.T, db *sqlx.DB, repo *RuleRepository) (channelIDNew string, ruleID string)
		rule      models.Rule
		wantRule  models.Rule
		wantErr   error
	}{
		"update all values": {
			setupData: setupData,
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{Class: "class2", Name: "read-only,ignored", ServiceID: "read-only,ignored"},
					},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "set below in test", Name: channel2.Name, Type: "read-only,ignored"},
					Recipient: "new@mail.com",
				},
				Active: true,
			},
			wantRule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin2",
							Class:     "class2",
							ServiceID: "service2",
						},
					},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "set below in test", Name: channel2.Name, Type: channel2.Type},
					Recipient: "new@mail.com",
				},
				Active: true,
			},
		},
		"update with non-existent origin should fail": {
			setupData: setupData,
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "non-existent"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test"},
				},
			},
			wantErr: ErrOriginsNotFound,
		},
		"update rule with non-existent channel works, but returns an empty channel ID": {
			setupData: func(t *testing.T, db *sqlx.DB, repo *RuleRepository) (channelID string, ruleID string) {
				_, ruleID = setupData(t, db, repo)
				return "", ruleID // return non-existent channel ID
			},
			rule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "class2"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: "set below in test"},
				},
			},
			wantRule: models.Rule{
				Name: "Updated Rule",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin2",
							Class:     "class2",
							ServiceID: "service2",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{}, // channel not found, so empty channel reference expected
				},
			},
		},
		"update with duplicate name should fail": {
			setupData: setupData,
			rule: models.Rule{
				Name: existingUntouchedRule.Name, // duplicate name
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "class2"}},
				},
				Action: models.Action{
					Channel:   models.ChannelReference{ID: "set below in test", Name: channel2.Name, Type: channel2.Type},
					Recipient: "",
				},
				Active: true,
			},
			wantErr: ErrDuplicateRuleName,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := pgtesting.NewDB(t)

			repo, err := NewRuleRepository(db)
			require.NoError(t, err)

			ctx := context.Background()

			channelIDNew, ruleID := tt.setupData(t, db, repo)
			tt.wantRule.ID = ruleID
			if tt.wantRule.Action.Channel.ID != "" {
				tt.wantRule.Action.Channel.ID = channelIDNew
			}
			tt.rule.ID = ruleID                      // set ID for update
			tt.rule.Action.Channel.ID = channelIDNew // set channel ID for update

			updatedRule, err := repo.Update(ctx, ruleID, tt.rule)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantRule, updatedRule)
			// Retrieve rule and verify
			gotRule, err := repo.Get(ctx, ruleID)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRule, gotRule)
		})
	}
}

func Test_ListRules(t *testing.T) {
	t.Parallel()

	t.Run("get empty list of rules", func(t *testing.T) {
		t.Parallel()
		db := pgtesting.NewDB(t)
		repo, err := NewRuleRepository(db)
		require.NoError(t, err)

		ctx := context.Background()

		rules, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, rules)
	})

	t.Run("get empty list of rules", func(t *testing.T) {
		t.Parallel()
		db := pgtesting.NewDB(t)
		repo, err := NewRuleRepository(db)
		require.NoError(t, err)

		ctx := context.Background()

		rules, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, rules)
	})

	t.Run("get all rules ordered by name and references to non-existent channel are returned empty", func(t *testing.T) {
		t.Parallel()
		db := pgtesting.NewDB(t)
		repo, err := NewRuleRepository(db)
		require.NoError(t, err)
		channelID1 := createTestChannel(t, db, "test-channel1", "mattermost")
		channelID2 := createTestChannel(t, db, "test-channel2", "teams")
		createTestOrigin(t, db, "Origin1", "class1", "service1")
		createTestOrigin(t, db, "Origin2", "class2", "service2")

		ctx := context.Background()

		rulesIn := []models.Rule{
			{
				Name: "Rule 2",
				Trigger: models.Trigger{
					Levels:  []string{"low"},
					Origins: []models.OriginReference{{Class: "class2"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: channelID2},
				},
				Active: true,
			},
			{
				Name: "Rule 1",
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "class1"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: channelID1},
				},
				Active: true,
			},
			{
				Name: "Rule 3",
				Trigger: models.Trigger{
					Levels:  []string{"high"},
					Origins: []models.OriginReference{{Class: "class1"}},
				},
				Action: models.Action{
					Channel: models.ChannelReference{ID: uuid.NewString()}, // non-existent channel ID
				},
				Active: true,
			},
		}
		wantRules := []models.Rule{
			{
				Name: "Rule 1",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin1",
							Class:     "class1",
							ServiceID: "service1",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{
						ID:   channelID1,
						Name: "test-channel1",
						Type: "mattermost",
					},
				},
				Active: true,
			},
			{
				Name: "Rule 2",
				Trigger: models.Trigger{
					Levels: []string{"low"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin2",
							Class:     "class2",
							ServiceID: "service2",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{
						ID:   channelID2,
						Name: "test-channel2",
						Type: "teams",
					},
				},
				Active: true,
			},
			{
				Name: "Rule 3",
				Trigger: models.Trigger{
					Levels: []string{"high"},
					Origins: []models.OriginReference{
						{
							Name:      "Origin1",
							Class:     "class1",
							ServiceID: "service1",
						},
					},
				},
				Action: models.Action{
					Channel: models.ChannelReference{}, // channel not found, so empty channel reference expected
				},
				Active: true,
			},
		}
		require.Equal(t, len(rulesIn), len(wantRules), "test setup error: rulesIn and wantRules must have same length")

		for i := range rulesIn {
			_, err := repo.Create(ctx, rulesIn[i])
			require.NoError(t, err)
		}

		gotRules, err := repo.List(ctx)
		for i := range gotRules {
			require.NotEmpty(t, gotRules[i].ID)
			gotRules[i].ID = "" // set ID to empty for comparison as they are not known beforehand
		}
		require.NoError(t, err)
		assert.Equal(t, wantRules, gotRules)
	})
}

func Test_DeleteRule(t *testing.T) {
	t.Parallel()
	db := pgtesting.NewDB(t)
	repo, err := NewRuleRepository(db)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("deleting non-existing rule is a no-op", func(t *testing.T) {
		t.Parallel()
		err := repo.Delete(ctx, uuid.NewString())
		assert.NoError(t, err)
	})

	t.Run("insert and delete rule", func(t *testing.T) {
		t.Parallel()
		channelID := createTestChannel(t, db, "test-channel", "mattermost")
		createTestOrigin(t, db, "Test Origin", "class1", "test-ns")
		rule := models.Rule{
			Name: "Test Rule 11",
			Trigger: models.Trigger{
				Levels:  []string{"high"},
				Origins: []models.OriginReference{{Class: "class1"}},
			},
			Action: models.Action{
				Channel: models.ChannelReference{ID: channelID},
			},
			Active: true,
		}

		rule, err := repo.Create(ctx, rule)
		require.NoError(t, err)

		_, err = repo.Get(ctx, rule.ID) // verify source exist
		require.NoError(t, err)

		err = repo.Delete(ctx, rule.ID)
		assert.NoError(t, err)

		_, err = repo.Get(ctx, rule.ID)
		assert.ErrorIs(t, err, errs.ErrItemNotFound)
	})
}
