// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationservice

import (
	"context"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/notifications"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NotificationService_CreateNotification_Failure(t *testing.T) {

	// received notification
	notification := models.Notification{
		Origin:      "Test Origin",
		OriginClass: "/serviceID/origin1",
		Timestamp:   "2024-01-01T00:00:00Z",
		Title:       "Test Notification",
		Detail:      "This is a test notification",
		Level:       notifications.LevelInfo,
	}

	synctest.Test(t, func(t *testing.T) {
		mockNotificationRepo := mocks.NewNotificationRepository(t)
		ruleService := mocks.NewRuleService(t) // no config, any call to the rule service in this test would be wrong behavior

		// setup mock
		mockNotificationRepo.EXPECT().CreateNotification(mock.Anything, notification).Return(notification, assert.AnError).Once()

		// no config of further mocks, as they are not expected to be called in this test
		notificationService := NewNotificationService(
			mockNotificationRepo, ruleService, nil, nil, nil, nil).(*notificationService)

		defer notificationService.cancelForwardRetriesWorker()

		_, err := notificationService.CreateNotification(context.Background(), notification)
		require.Error(t, err)

		synctest.Wait()
	})

}

func Test_NotificationService_CreateNotification_Forwarding(t *testing.T) {
	// Test verifies that notifications can be forwarded to all channel types.
	// Additionally it verifies that failure in forwarding to one recipient/channel
	// does not affect other recipients/channels.

	synctest.Test(t, func(t *testing.T) {
		notification := models.Notification{
			Origin:      "Test Origin",
			OriginClass: "/serviceID/origin1",
			Timestamp:   "2024-01-01T00:00:00Z",
			Title:       "Test Notification",
			Detail:      "This is a test notification",
			Level:       notifications.LevelInfo,
		}

		matchMailSubject := mock.MatchedBy(func(subject string) bool {
			return strings.Contains(subject, notification.Title)
		})
		matchMessage := mock.MatchedBy(func(message string) bool {
			return strings.Contains(message, notification.Title) && strings.Contains(message, notification.Detail)
		})

		// Create three channels, each with a different type
		mattermostChannel := models.NotificationChannel{
			Id:          "mattermost-channel-id",
			ChannelType: models.ChannelTypeMattermost,
			ChannelName: "Mattermost Channel",
			WebhookUrl:  new("https://mattermost.example.com/webhook"),
		}

		teamsChannel := models.NotificationChannel{
			Id:          "teams-channel-id",
			ChannelType: models.ChannelTypeTeams,
			ChannelName: "Teams Channel",
			WebhookUrl:  new("https://teams.example.com/webhook"),
		}

		mailChannel := models.NotificationChannel{
			Id:          "mail-channel-id",
			ChannelType: models.ChannelTypeMail,
			ChannelName: "Mail Channel",
		}

		actions := []models.Action{ // simulates three matching rules
			{
				Channel: models.ChannelReference{
					ID:   mattermostChannel.Id,
					Type: mattermostChannel.ChannelType,
				},
			},
			{
				Channel: models.ChannelReference{
					ID:   teamsChannel.Id,
					Type: teamsChannel.ChannelType,
				},
			},
			{
				Channel: models.ChannelReference{
					ID:   mailChannel.Id,
					Type: mailChannel.ChannelType,
				},
				Recipient: "a@example.com, b@example.com",
			},
		}

		// Setup mocks
		mockNotificationRepo := mocks.NewNotificationRepository(t)
		ruleService := mocks.NewRuleService(t)
		channelService := mocks.NewNotificationChannelService(t)
		mailService := mocks.NewMailService(t)
		mattermostService := mocks.NewWebhookService(t)
		teamsService := mocks.NewWebhookService(t)

		mockNotificationRepo.EXPECT().CreateNotification(mock.Anything, notification).Return(notification, nil).Once()
		ruleService.EXPECT().ProcessRules(mock.Anything, notification).Return(actions, nil)

		// Mock channel service calls - all three channels should be fetched
		channelService.EXPECT().GetNotificationChannelByIdAndType(mock.Anything, mattermostChannel.Id, mattermostChannel.ChannelType).
			Return(mattermostChannel, nil).Once()
		channelService.EXPECT().GetNotificationChannelByIdAndType(mock.Anything, teamsChannel.Id, teamsChannel.ChannelType).
			Return(teamsChannel, nil).Once()
		channelService.EXPECT().GetNotificationChannelByIdAndType(mock.Anything, mailChannel.Id, mailChannel.ChannelType).
			Return(mailChannel, nil).Once()

		// Mock forwarding services
		// Rule/Action 1 (Mattermost) - should succeed
		mattermostService.EXPECT().SendMessage(
			*mattermostChannel.WebhookUrl,
			matchMessage,
		).Return(nil).Once()

		// Rule/Action 2 (Teams) - fails sending
		teamsService.EXPECT().SendMessage(
			*teamsChannel.WebhookUrl,
			matchMessage,
		).Return(assert.AnError).Once()

		// Rule/Action 3 (Mail) - should send despite rule 2 failing
		// and failure on first recipient should not affect second recipient
		mailService.EXPECT().SendMail(
			mock.Anything,
			mailChannel,
			"a@example.com",
			matchMailSubject,
			notification.Detail,
		).Return(assert.AnError).Once()
		mailService.EXPECT().SendMail(
			mock.Anything,
			mailChannel,
			"b@example.com",
			matchMailSubject,
			notification.Detail,
		).Return(nil).Once()

		notificationService := NewNotificationService(
			mockNotificationRepo,
			ruleService,
			channelService,
			mailService,
			mattermostService,
			teamsService,
		).(*notificationService)

		defer notificationService.cancelForwardRetriesWorker()

		// Create notification should succeed even though some forwarding attempts fail
		_, err := notificationService.CreateNotification(context.Background(), notification)
		require.NoError(t, err)

		synctest.Wait()
	})
}

func Test_NotificationService_RetryLogic_MaxRetriesReached(t *testing.T) {
	t.Parallel()

	notification := models.Notification{
		Origin:      "Test Origin",
		OriginClass: "/serviceID/origin1",
		Timestamp:   "2024-01-01T00:00:00Z",
		Title:       "Test Notification",
		Detail:      "This is a test notification",
		Level:       notifications.LevelInfo,
	}

	matchMailSubject := mock.MatchedBy(func(subject string) bool {
		return strings.Contains(subject, notification.Title)
	})
	matchMessage := mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, notification.Title) && strings.Contains(message, notification.Detail)
	})

	mailchannel := models.NotificationChannel{
		Id:          "mail-channel-id",
		ChannelType: models.ChannelTypeMail,
		ChannelName: "Mail Channel",
	}
	teamsChannel := models.NotificationChannel{
		Id:          "teams-channel-id",
		ChannelType: models.ChannelTypeTeams,
		ChannelName: "Teams Channel",
		WebhookUrl:  new("https://teams.example.com/webhook"),
	}
	mattermostChannel := models.NotificationChannel{
		Id:          "mattermost-channel-id",
		ChannelType: models.ChannelTypeMattermost,
		ChannelName: "Mattermost Channel",
		WebhookUrl:  new("https://mattermost.example.com/webhook"),
	}

	tests := map[string]struct {
		processRuleFailures int
		action              models.Action
		mockConfig          func(
			t *testing.T,
			channelService *mocks.NotificationChannelService,
			mailService *mocks.MailService,
			mattermostService, teamsService *mocks.WebhookService,
		)
	}{
		"ProcessRules is retried up to max retries": {
			processRuleFailures: maxRetries + 1,
			// no more mock calls, as operation is aborted
		},
		"Getting channel is retried up to max retries": {
			action: models.Action{
				Channel: models.ChannelReference{
					ID:   teamsChannel.Id,
					Type: models.ChannelTypeTeams,
				},
			},
			mockConfig: func(t *testing.T, channelService *mocks.NotificationChannelService, _ *mocks.MailService, _, _ *mocks.WebhookService) {
				channelService.EXPECT().GetNotificationChannelByIdAndType(
					mock.Anything,
					teamsChannel.Id,
					models.ChannelTypeTeams,
				).Return(models.NotificationChannel{}, assert.AnError).Times(maxRetries + 1)
			},
		},
		"Mail send is retried up to max retries": {
			action: models.Action{
				Channel: models.ChannelReference{
					ID:   mailchannel.Id,
					Type: mailchannel.ChannelType,
				},
				Recipient: "success@example.com,failure@example.com,maxRetries@example.com",
			},
			mockConfig: func(t *testing.T, channelService *mocks.NotificationChannelService, mailService *mocks.MailService, _, _ *mocks.WebhookService) {
				channelService.EXPECT().GetNotificationChannelByIdAndType(mock.Anything, mailchannel.Id, models.ChannelTypeMail).
					Return(mailchannel, nil).Times(maxRetries + 1 + maxRetries) // initial attempt needs only a single fetch for all recipients

				mailService.EXPECT().SendMail(
					mock.Anything,
					mailchannel,
					"success@example.com",
					matchMailSubject,
					notification.Detail,
				).Return(nil).Once()

				mailService.EXPECT().SendMail(
					mock.Anything,
					mailchannel,
					"failure@example.com",
					mock.MatchedBy(func(subject string) bool {
						return strings.Contains(subject, notification.Title)
					}),
					notification.Detail,
				).Return(assert.AnError).Times(maxRetries + 1)

				mailService.EXPECT().SendMail(
					mock.Anything,
					mailchannel,
					"maxRetries@example.com",
					matchMailSubject,
					notification.Detail,
				).Return(assert.AnError).Times(maxRetries)

				mailService.EXPECT().SendMail(
					mock.Anything,
					mailchannel,
					"maxRetries@example.com",
					matchMailSubject,
					notification.Detail,
				).Return(nil).Once()
			},
		},
		"Mattermost send is retried up to max retries": {
			action: models.Action{
				Channel: models.ChannelReference{
					ID:   mattermostChannel.Id,
					Type: models.ChannelTypeMattermost,
				},
			},
			mockConfig: func(t *testing.T, channelService *mocks.NotificationChannelService, _ *mocks.MailService, mattermostService, _ *mocks.WebhookService) {
				channelService.EXPECT().GetNotificationChannelByIdAndType(
					mock.Anything,
					mattermostChannel.Id,
					models.ChannelTypeMattermost,
				).Return(mattermostChannel, nil).Times(maxRetries + 1)

				mattermostService.EXPECT().SendMessage(
					*mattermostChannel.WebhookUrl,
					matchMessage,
				).Return(assert.AnError).Times(maxRetries + 1)
			},
		},
		"Teams send is retried up to max retries": {
			action: models.Action{
				Channel: models.ChannelReference{
					ID:   teamsChannel.Id,
					Type: models.ChannelTypeTeams,
				},
			},
			mockConfig: func(t *testing.T, channelService *mocks.NotificationChannelService, _ *mocks.MailService, _, teamsService *mocks.WebhookService) {
				channelService.EXPECT().GetNotificationChannelByIdAndType(
					mock.Anything,
					teamsChannel.Id,
					models.ChannelTypeTeams,
				).Return(teamsChannel, nil).Times(maxRetries + 1)

				teamsService.EXPECT().SendMessage(
					*teamsChannel.WebhookUrl,
					matchMessage,
				).Return(assert.AnError).Times(maxRetries + 1)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				mockNotificationRepo := mocks.NewNotificationRepository(t)
				ruleService := mocks.NewRuleService(t)
				channelService := mocks.NewNotificationChannelService(t)
				mailService := mocks.NewMailService(t)
				mattermostService := mocks.NewWebhookService(t)
				teamsService := mocks.NewWebhookService(t)

				notificationService := NewNotificationService(
					mockNotificationRepo,
					ruleService,
					channelService,
					mailService,
					mattermostService,
					teamsService,
				).(*notificationService)

				// stop the worker to avoid go routines leak
				defer notificationService.cancelForwardRetriesWorker()

				mockNotificationRepo.EXPECT().CreateNotification(mock.Anything, notification).Return(notification, nil).Once()

				if tt.processRuleFailures > 0 {
					ruleService.EXPECT().ProcessRules(mock.Anything, notification).Return(nil, assert.AnError).Times(tt.processRuleFailures)
				}
				if tt.processRuleFailures <= maxRetries {
					ruleService.EXPECT().ProcessRules(mock.Anything, notification).Return([]models.Action{tt.action}, nil).Once()
				}

				if tt.mockConfig != nil {
					tt.mockConfig(t, channelService, mailService, mattermostService, teamsService)
				}

				// Create notification - this triggers initial send which will fail
				_, err := notificationService.CreateNotification(context.Background(), notification)
				require.NoError(t, err)

				// wait until all retries have been processed
				// note: use generous duration, failing with exponential backoff after max retries
				// takes roughly `baseDelayRetryForwarding*(2^(maxRetries+1)-1)`
				time.Sleep(baseDelayRetryForwarding*(1<<uint(maxRetries+1)-1) + 5*time.Hour)
				synctest.Wait()
			})
		})
	}

}

func Test_NotificationService_RetryLogic_QueueOverflow(t *testing.T) {
	// Test verifies that at maximum `maxRetainedFailedSends` failed send attempts are retained for retrying.
	// Failed sends beyond the capacity should be dropped gracefully and not block the service.

	notification := models.Notification{
		Origin:      "Test Origin",
		OriginClass: "/serviceID/origin1",
		Timestamp:   "2024-01-01T00:00:00Z",
		Title:       "Test Notification",
		Detail:      "This is a test notification",
		Level:       notifications.LevelInfo,
	}

	teamsChannel := models.NotificationChannel{
		Id:          "teams-channel-id",
		ChannelType: models.ChannelTypeTeams,
		ChannelName: "Teams Channel",
		WebhookUrl:  new("https://teams.example.com/webhook"),
	}

	// create more actions than the queue can hold
	actions := make([]models.Action, maxRetainedFailedSends+10)
	for i := range actions {
		actions[i] = models.Action{
			Channel: models.ChannelReference{
				ID:   teamsChannel.Id,
				Type: teamsChannel.ChannelType,
			},
		}
	}

	synctest.Test(t, func(t *testing.T) {
		mockNotificationRepo := mocks.NewNotificationRepository(t)
		ruleService := mocks.NewRuleService(t)
		channelService := mocks.NewNotificationChannelService(t)
		teamsService := mocks.NewWebhookService(t)

		mockNotificationRepo.EXPECT().CreateNotification(mock.Anything, notification).Return(notification, nil).Once()
		ruleService.EXPECT().ProcessRules(mock.Anything, notification).Return(actions, nil)

		channelService.EXPECT().GetNotificationChannelByIdAndType(
			mock.Anything,
			teamsChannel.Id,
			teamsChannel.ChannelType,
		).Return(teamsChannel, nil).
			Times(len(actions) + maxRetainedFailedSends) // initial attempts + retries for the ones that are not dropped

		// Mock SendMessage to fail every time for all initial attempts
		teamsService.EXPECT().SendMessage(
			*teamsChannel.WebhookUrl,
			mock.MatchedBy(func(message string) bool {
				return strings.Contains(message, notification.Title)
			}),
		).Return(assert.AnError).Times(len(actions))
		// only the first maxRetainedFailedSends are exptected to be retried
		teamsService.EXPECT().SendMessage(
			*teamsChannel.WebhookUrl,
			mock.MatchedBy(func(message string) bool {
				return strings.Contains(message, notification.Title)
			}),
		).Return(nil).Times(maxRetainedFailedSends)

		notificationService := NewNotificationService(
			mockNotificationRepo,
			ruleService,
			channelService,
			nil,
			nil,
			teamsService,
		).(*notificationService)

		defer notificationService.cancelForwardRetriesWorker()

		// Create notification - this triggers 250 failed sends
		_, err := notificationService.CreateNotification(context.Background(), notification)
		require.NoError(t, err, "failures in forwarding should not fail notification creation")

		// wait until all retries have been processed
		// note: use generous duration, failing with exponential backoff after one retry
		// takes roughly `baseDelayRetryForwarding`
		time.Sleep(baseDelayRetryForwarding * 10)
		synctest.Wait()
	})
}
