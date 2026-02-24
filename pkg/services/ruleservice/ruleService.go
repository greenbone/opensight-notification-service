// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package ruleservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/greenbone/opensight-notification-service/pkg/errs"
	"github.com/greenbone/opensight-notification-service/pkg/models"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
)

var ErrRuleLimitReached = fmt.Errorf("alert rule limit reached")
var ErrRecipientRequired = fmt.Errorf("recipient is required for the selected channel")
var ErrRecipientNotSupported = fmt.Errorf("recipient is not supported for the selected channel")
var ErrChannelNotFound = fmt.Errorf("notification channel not found")

type RuleRepository interface {
	Get(ctx context.Context, id string) (models.Rule, error)
	List(ctx context.Context) ([]models.Rule, error)
	Create(ctx context.Context, rule models.Rule) (models.Rule, error)
	Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error)
	Delete(ctx context.Context, id string) error
}

type RuleService struct {
	store        RuleRepository
	channelStore notificationrepository.NotificationChannelRepository
	ruleLimit    int
}

func NewRuleService(store RuleRepository, channelStore notificationrepository.NotificationChannelRepository, ruleLimit int) *RuleService {
	return &RuleService{
		store:        store,
		channelStore: channelStore,
		ruleLimit:    ruleLimit,
	}
}

func (s *RuleService) Get(ctx context.Context, id string) (models.Rule, error) {
	return s.store.Get(ctx, id)
}

func (s *RuleService) List(ctx context.Context) ([]models.Rule, error) {
	return s.store.List(ctx)
}

func (s *RuleService) Create(ctx context.Context, rule models.Rule) (models.Rule, error) {
	rules, err := s.store.List(ctx)
	if err != nil {
		return models.Rule{}, fmt.Errorf("failed to check the rule limit")
	}
	if len(rules) >= s.ruleLimit {
		return models.Rule{}, ErrRuleLimitReached
	}

	err = s.validateAction(ctx, rule.Action)
	if err != nil {
		return models.Rule{}, err
	}

	return s.store.Create(ctx, rule)
}

func (s *RuleService) Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error) {
	err := s.validateAction(ctx, rule.Action)
	if err != nil {
		return models.Rule{}, err
	}

	return s.store.Update(ctx, id, rule)
}

func (s *RuleService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *RuleService) validateAction(ctx context.Context, action models.Action) error {
	channel, err := s.channelStore.GetNotificationChannelById(ctx, action.Channel.ID)
	if err != nil {
		if errors.Is(err, errs.ErrItemNotFound) {
			return ErrChannelNotFound // from perspective of the service this is not a generic not found, but an issue with the passed object
		}
		return fmt.Errorf("failed to get notification channel: %w", err)
	}

	if channel.ChannelType == string(models.ChannelTypeMail) {
		if action.Recipient == "" {
			return ErrRecipientRequired
		}
	} else if action.Recipient != "" {
		return ErrRecipientNotSupported
	}

	return nil
}
