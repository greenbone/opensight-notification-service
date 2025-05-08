// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package ruleservice

import (
	"context"

	"github.com/greenbone/opensight-notification-service/pkg/models"
)

type RuleRepository interface {
	Get(ctx context.Context, id string) (models.Rule, error)
	List(ctx context.Context) ([]models.Rule, error)
	Create(ctx context.Context, rule models.Rule) (models.Rule, error)
	Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error)
	Delete(ctx context.Context, id string) error
}

type RuleService struct {
	store RuleRepository
}

func NewRuleService(store RuleRepository) *RuleService {
	return &RuleService{store: store}
}

func (s *RuleService) Get(ctx context.Context, id string) (models.Rule, error) {
	return s.store.Get(ctx, id)
}

func (s *RuleService) List(ctx context.Context) ([]models.Rule, error) {
	return s.store.List(ctx)
}

func (s *RuleService) Create(ctx context.Context, rule models.Rule) (models.Rule, error) {
	return s.store.Create(ctx, rule)
}

func (s *RuleService) Update(ctx context.Context, id string, rule models.Rule) (models.Rule, error) {
	return s.store.Update(ctx, id, rule)
}

func (s *RuleService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
