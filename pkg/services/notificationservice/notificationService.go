// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationservice

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/logs"
	"github.com/greenbone/opensight-golang-libraries/pkg/notifications"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-notification-service/pkg/models"
)

const (
	maxRetries                   = 10
	baseDelayRetryRuleProcessing = 10 * time.Second
	baseDelayRetryForwarding     = time.Minute
	retryPollInterval            = 5 * time.Second // intervall to check for pending send tasks
	maxRetainedFailedSends       = 400             // buffer size for failed sends, arbitrary limit to avoid memory issues
)

// multipleNewlinesRegex matches 2 or more consecutive newlines
var multipleNewlinesRegex = regexp.MustCompile(`\n{2,}`)

type NotificationService interface {
	ListNotifications(
		ctx context.Context,
		resultSelector query.ResultSelector,
	) (notifications []models.Notification, totalResult uint64, err error)
	CreateNotification(
		ctx context.Context,
		notificationIn models.Notification,
	) (notification models.Notification, err error)
}

type NotificationRepository interface {
	ListNotifications(
		ctx context.Context,
		resultSelector query.ResultSelector,
	) (notifications []models.Notification, totalResult uint64, err error)
	CreateNotification(ctx context.Context, notification models.Notification) (models.Notification, error)
}

type RuleService interface {
	ProcessRules(ctx context.Context, notification models.Notification) ([]models.Action, error)
}

type NotificationChannelService interface {
	GetNotificationChannelByIdAndType(
		ctx context.Context,
		id string,
		channelType models.ChannelType,
	) (models.NotificationChannel, error)
}

type WebhookService interface {
	SendMessage(webhookUrl string, message string) error
}

type MailService interface {
	SendMail(
		ctx context.Context,
		channel models.NotificationChannel,
		recipient string,
		subject string,
		htmlBody string,
	) error
}

type SendTask struct {
	ctx           context.Context
	Notification  *models.Notification // avoid copies, as object can be quite large
	Action        models.Action
	attempt       int
	nextExecution time.Time
}

type notificationService struct {
	store             NotificationRepository
	ruleService       RuleService
	channelService    NotificationChannelService
	mailService       MailService
	mattermostService WebhookService
	teamsService      WebhookService

	failedSends chan SendTask
	// only for tests: allows to shut down the forward retries worker to avoid goroutine leaks
	cancelForwardRetriesWorker context.CancelFunc
}

func NewNotificationService(
	store NotificationRepository,
	ruleService RuleService,
	channelService NotificationChannelService,
	mailService MailService,
	mattermostService WebhookService,
	teamsService WebhookService,
) NotificationService {

	service := &notificationService{
		store:             store,
		ruleService:       ruleService,
		channelService:    channelService,
		mailService:       mailService,
		mattermostService: mattermostService,
		teamsService:      teamsService,
		failedSends:       make(chan SendTask, maxRetainedFailedSends),
	}

	contextForwardRetriesWorker, cancel := context.WithCancel(context.Background())
	service.cancelForwardRetriesWorker = cancel
	go service.forwardRetriesWorker(contextForwardRetriesWorker)

	return service
}

func (s *notificationService) ListNotifications(
	ctx context.Context,
	resultSelector query.ResultSelector,
) (notifications []models.Notification, totalResult uint64, err error) {
	return s.store.ListNotifications(ctx, resultSelector)
}

func (s *notificationService) CreateNotification(
	ctx context.Context,
	notificationIn models.Notification,
) (models.Notification, error) {

	notification, err := s.store.CreateNotification(ctx, notificationIn)
	if err != nil {
		return models.Notification{}, fmt.Errorf("failed to store notification: %w", err)
	}

	// only process rules after notification was successfully stored, this avoids forwarding
	// the notification multiple times if the client retries creating notification
	go func() {
		ctxForward := context.WithoutCancel(ctx)
		attempt := 0
		for {
			err := s.processRules(ctxForward, notificationIn)
			if err == nil {
				break
			}
			logs.Ctx(ctx).Err(err).Int("attempt", attempt).Msg("failed to process rules")
			if attempt >= maxRetries {
				logs.Ctx(ctxForward).Error().Err(err).Int("attempt", attempt+1).Msg("Skip processing of rules after maximum of retries")
				break
			}
			time.Sleep(exponentialBackoff(baseDelayRetryRuleProcessing, attempt))
			attempt++
		}
	}()

	return notification, nil
}

func (s *notificationService) processRules(ctx context.Context, notification models.Notification) error {
	actions, err := s.ruleService.ProcessRules(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to process rules: %w", err)
	}

	for _, action := range actions {
		sendTask := SendTask{
			ctx:          context.WithoutCancel(ctx),
			Notification: &notification,
			Action:       action,
		}
		s.forwardNotification(sendTask)
	}

	return nil
}

// forwardNotification sends the notification according to the action.
// If sending fails, it is scheduled for retry with exponential backoff.
func (s *notificationService) forwardNotification(sendTask SendTask) {
	ctx := sendTask.ctx
	notification := *sendTask.Notification
	action := sendTask.Action

	subject := createSubject(notification)
	body := notification.Detail

	channel, err := s.channelService.GetNotificationChannelByIdAndType(ctx, action.Channel.ID, action.Channel.Type)
	if err != nil {
		logs.Ctx(ctx).Err(err).Int("attempt", sendTask.attempt).Msg("failed to get channel for forwarding notification")
		s.scheduleRetry(sendTask)
		return
	}

	switch channelType := action.Channel.Type; channelType {
	case models.ChannelTypeMail:
		err := s.mailService.SendMail(ctx, channel, action.Recipient, subject, body)
		if err != nil {
			logs.Ctx(ctx).Err(err).Int("attempt", sendTask.attempt).Msg("failed to send mail")
			s.scheduleRetry(sendTask)
		}
	case models.ChannelTypeTeams:
		err := s.teamsService.SendMessage(*channel.WebhookUrl, convertToMarkDownMessage(subject, body))
		if err != nil {
			logs.Ctx(ctx).Err(err).Int("attempt", sendTask.attempt).Msg("failed to send teams message")
			s.scheduleRetry(sendTask)
		}

	case models.ChannelTypeMattermost:
		err := s.mattermostService.SendMessage(*channel.WebhookUrl, convertToMarkDownMessage(subject, body))
		if err != nil {
			logs.Ctx(ctx).Err(err).Int("attempt", sendTask.attempt).Msg("failed to send mattermost message")
			s.scheduleRetry(sendTask)
		}

	default:
		logs.Ctx(ctx).Error().Msgf("invalid channel type: %s, allowed are %v", channelType, models.AllowedChannels)
	}
}

// scheduleRetry calculates the next execution time using exponential backoff and queues the task for retry.
// If the queue is full or the maximum number of retries has been reached, the message will be dropped.
func (s *notificationService) scheduleRetry(sendTask SendTask) {
	if sendTask.attempt >= maxRetries {
		logs.Ctx(sendTask.ctx).Error().
			Str("channel", sendTask.Action.Channel.ID).
			Str("channelName", sendTask.Action.Channel.Name).
			Str("channelType", string(sendTask.Action.Channel.Type)).
			Int("retries", sendTask.attempt).
			Msg("Dropping message after maximum of retries")
		return
	}
	sendTask.nextExecution = time.Now().Add(exponentialBackoff(baseDelayRetryForwarding, sendTask.attempt))
	sendTask.attempt++

	select {
	case s.failedSends <- sendTask:
	default:
		logs.Ctx(sendTask.ctx).Error().
			Str("channel", sendTask.Action.Channel.ID).
			Str("channelName", sendTask.Action.Channel.Name).
			Str("channelType", string(sendTask.Action.Channel.Type)).
			Msg("Dropping message, retry queue is full")
	}
}

// exponentialBackoff calculates the backoff, retry count is 0-indexed.
func exponentialBackoff(baseDelay time.Duration, retryCount int) time.Duration {
	backoff := baseDelay * (1 << retryCount)                                 // Exponential backoff: baseDelay * 2^retryCount
	jitter := time.Duration(float64(backoff) * 0.2 * (0.5 - rand.Float64())) // Add jitter: +/-10% of the backoff
	return backoff + jitter
}

// forwardRetriesWorker continuously listens for failed send tasks and retries them when their next execution time has come.
// Send attempts are only done periodically to save cpu load.
func (s *notificationService) forwardRetriesWorker(ctx context.Context) {
	pendingSendTasks := make([]SendTask, 0, maxRetainedFailedSends)

	tick := time.Tick(retryPollInterval)
	for {
		select {
		case sendTask := <-s.failedSends:
			if len(pendingSendTasks) >= maxRetainedFailedSends {
				logs.Ctx(sendTask.ctx).Error().
					Str("channel", sendTask.Action.Channel.ID).
					Str("channelName", sendTask.Action.Channel.Name).
					Str("channelType", string(sendTask.Action.Channel.Type)).
					Msg("Dropping message, retry queue is full")
				continue
			}
			pendingSendTasks = append(pendingSendTasks, sendTask)
		case <-tick:
			newIndex := 0
			for i := range pendingSendTasks {
				if time.Now().Before(pendingSendTasks[i].nextExecution) {
					pendingSendTasks[newIndex] = pendingSendTasks[i]
					newIndex++
					continue
				}
				s.forwardNotification(pendingSendTasks[i])
			}
			pendingSendTasks = pendingSendTasks[:newIndex] // remove tasks that were just sent
		case <-ctx.Done():
			return
		}
	}
}

func createSubject(notification models.Notification) string {
	var icon string
	switch notification.Level {
	case notifications.LevelUrgent, notifications.LevelError:
		icon = "🔴"
	case notifications.LevelWarning:
		icon = "🟡"
	case notifications.LevelInfo:
		icon = "🔵"
	}

	subject := fmt.Sprintf("%s [%s]", notification.Title, notification.Origin)
	if icon != "" {
		subject = fmt.Sprintf("%s %s", icon, subject)
	}

	return subject
}

// convertToMarkDownMessage formats subject and body into valid markdown
// which is displayed well for both mattermost and teams
func convertToMarkDownMessage(subject, body string) string {
	const escapedNewlinePlaceholder = "\x00ESCAPED_NEWLINE\x00"
	const consecutiveNewlinesPlaceholder = "\x00CONSECUTIVE_NEWLINES\x00"

	// Temporarily replace escaped newlines (\\n) with a placeholder
	body = strings.ReplaceAll(body, `\\n`, escapedNewlinePlaceholder)
	// Convert unescaped literal \n strings to actual newlines
	body = strings.ReplaceAll(body, `\n`, "\n")
	// Protect all sequences of 2+ consecutive newlines by replacing each with a placeholder
	body = multipleNewlinesRegex.ReplaceAllStringFunc(body, func(match string) string {
		return strings.Repeat(consecutiveNewlinesPlaceholder, len(match))
	})
	// Now all remaining single newlines become double newlines
	body = strings.ReplaceAll(body, "\n", "\n\n")
	// Restore consecutive newlines (each placeholder back to a newline)
	body = strings.ReplaceAll(body, consecutiveNewlinesPlaceholder, "\n")
	// Restore escaped newlines as literal \n text
	body = strings.ReplaceAll(body, escapedNewlinePlaceholder, `\n`)

	markdownMessage := fmt.Sprintf("**%s**\n\n%s", subject, body)

	return markdownMessage
}
