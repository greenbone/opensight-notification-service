// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originservice

import (
	"context"

	"github.com/greenbone/opensight-notification-service/pkg/entities"
)

type OriginRepository interface {
	UpsertOrigins(ctx context.Context, serviceID string, origins []entities.Origin) error
	ListOrigins(ctx context.Context) ([]entities.Origin, error)
}

type OriginService struct {
	store OriginRepository
}

func NewOriginService(store OriginRepository) *OriginService {
	return &OriginService{store: store}
}

func (s *OriginService) UpsertOrigins(ctx context.Context, serviceID string, origins []entities.Origin) error {
	return s.store.UpsertOrigins(ctx, serviceID, origins)
}

func (s *OriginService) ListOrigins(ctx context.Context) ([]entities.Origin, error) {
	return s.store.ListOrigins(ctx)
}
