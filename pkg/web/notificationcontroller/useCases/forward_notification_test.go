// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package usecases

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-golang-libraries/pkg/httpassert"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/originrepository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/rulerepository"
	"github.com/greenbone/opensight-notification-service/pkg/security"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice/mocks"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/ruleservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/iam"
	"github.com/greenbone/opensight-notification-service/pkg/web/integrationTests"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/notificationcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/rulecontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/testhelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (
	*gin.Engine,
	*mocks.MailService,
) {
	// setup repositories
	db := pgtesting.NewDB(t)
	encryptMgr := security.NewEncryptManager()
	encryptMgr.UpdateKeys(config.DatabaseEncryptionKey{
		Password:     "password",
		PasswordSalt: "password-salt-should-no-be-short-fyi",
	})
	notificationRepo, err := notificationrepository.NewNotificationRepository(db)
	require.NoError(t, err)
	channelRepo, err := notificationrepository.NewNotificationChannelRepository(db, encryptMgr)
	require.NoError(t, err)
	ruleRepo, err := rulerepository.NewRuleRepository(db)
	require.NoError(t, err)
	originRepo, err := originrepository.NewOriginRepository(db)
	require.NoError(t, err)

	// setup services
	mockMailService := mocks.NewMailService(t)
	ruleLimit := 100
	channelService := notificationchannelservice.NewNotificationChannelService(channelRepo)
	mailLimit := 10
	mailChannelService := notificationchannelservice.NewMailChannelService(channelService, mockMailService, mailLimit)
	ruleService, err := ruleservice.NewRuleService(ruleRepo, channelRepo, originRepo, ruleLimit)
	require.NoError(t, err)

	notificationSvc := notificationservice.NewNotificationService(
		notificationRepo,
		ruleService,
		channelService,
		mockMailService,
		nil,
		nil,
	)

	registry := errmap.NewRegistry()
	router := testhelper.NewTestWebEngine(registry)
	authMiddleware, err := auth.NewGinAuthMiddleware(integrationTests.NewTestJwtParser(t))
	require.NoError(t, err)

	mailcontroller.NewMailController(router, channelService, mailChannelService, authMiddleware, registry)
	rulecontroller.NewRuleController(router, ruleService, authMiddleware, registry)
	notificationcontroller.AddNotificationController(router, notificationSvc, authMiddleware)

	return router, mockMailService
}

func TestForwardNotification(t *testing.T) {
	router, mockMailService := setup(t)

	// notification that should trigger a rule
	notification := models.Notification{
		Timestamp:   time.Now().Format(time.RFC3339Nano),
		Origin:      "Test Origin",
		OriginClass: "test-service/test-task",
		Title:       "Test title",
		Detail:      "This is a test notification that should be forwarded via email",
		Level:       "info",
	}

	// create mail channel
	mailChannel := testhelper.GetValidMailNotificationChannel()
	var mailChannelID string
	httpassert.New(t, router).
		Post("/notification-channel/mail").
		JsonContentObject(mailChannel).
		AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$.id", httpassert.ExtractTo(&mailChannelID))
	require.NotEmpty(t, mailChannelID)

	// create rule
	var ruleID string
	httpassert.New(t, router).
		Post("/rules").
		JsonContent(fmt.Sprintf(`
			{
				"name": "Test Forwarding Rule",
				"trigger": {
					"levels": ["info", "warning", "error", "urgent"],
					"origins": [{
						"class": "%s"
					}]
				},
				"action": {
					"channel": {
						"id": "%s",
						"type": "mail"
					},
					"recipient": " a@example.com ,   b@example.com "
				},
				"active": true
			}`, models.OriginAllClass, mailChannelID),
		).
		AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Admin)).
		Expect().
		StatusCode(http.StatusCreated).
		JsonPath("$.id", httpassert.ExtractTo(&ruleID))
	require.NotEmpty(t, ruleID)

	expectedRecipients := []string{"a@example.com", "b@example.com"}

	// configure mock
	notificationReceived := make(chan string, len(expectedRecipients))
	expectMessageFromRecipient := func(recipient string) {
		mockMailService.EXPECT().SendMail(
			mock.Anything,
			mock.Anything,
			recipient,
			mock.MatchedBy(func(subject string) bool {
				return strings.Contains(subject, notification.Title)
			}),
			mock.MatchedBy(func(body string) bool {
				return strings.Contains(body, notification.Detail)
			}),
		).RunAndReturn(func(ctx context.Context, channel models.NotificationChannel, recipient, subject, htmlBody string) error {
			notificationReceived <- recipient
			return nil
		}).Times(1)
	}
	for _, recipient := range expectedRecipients {
		expectMessageFromRecipient(recipient)
	}

	// create notification that should trigger the rule and be forwarded to all recipients
	httpassert.New(t, router).
		Post("/notifications").
		JsonContentObject(notification).
		AuthJwt(integrationTests.CreateJwtTokenWithRole(iam.Notification)).
		Expect().
		StatusCode(http.StatusCreated)

	// Wait for both notifications to be forwarded or timeout
	var receivedRecipients []string
	timeout := time.After(2 * time.Second)
	for range len(expectedRecipients) {
		select {
		case recipient := <-notificationReceived:
			receivedRecipients = append(receivedRecipients, recipient)
		case <-timeout:
			t.Fatalf("Timeout waiting for notification to be forwarded. Received %d/%d notifications to: %v",
				len(receivedRecipients),
				len(expectedRecipients),
				receivedRecipients,
			)
		}
	}

	require.ElementsMatch(t, expectedRecipients, receivedRecipients)
}
