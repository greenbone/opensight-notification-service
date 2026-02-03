// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpsertOrigins_ListOrigins(t *testing.T) {
	type input struct {
		origins   []entities.Origin
		serviceID string
	}

	tests := map[string]struct {
		inputs      []input
		wantOrigins []entities.Origin
		wantErr     bool
	}{
		"create origins from single service (no prior data)": {
			inputs: []input{
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA", ServiceID: "read only, to be ignored"},
						{Name: "origin2", Class: "classB"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin1", Class: "classA", ServiceID: "service1"},
				{Name: "origin2", Class: "classB", ServiceID: "service1"},
			},
		},
		"create origins from multiple services": {
			inputs: []input{
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
				{
					serviceID: "service2",
					origins: []entities.Origin{
						{Name: "origin2", Class: "classB"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin1", Class: "classA", ServiceID: "service1"},
				{Name: "origin2", Class: "classB", ServiceID: "service2"},
			},
		},
		"upsert origins replaces the entries from the same service": {
			inputs: []input{
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
						{Name: "origin2", Class: "classB"},
					},
				},
				{
					serviceID: "service2",
					origins: []entities.Origin{
						{Name: "origin3", Class: "classC"},
					},
				},
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin4", Class: "classD"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin4", Class: "classD", ServiceID: "service1"},
				{Name: "origin3", Class: "classC", ServiceID: "service2"},
			},
		},
		"delete all origins from a service when upserting empty list": {
			inputs: []input{
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
				{
					serviceID: "service1",
					origins:   []entities.Origin{},
				},
			},
			wantOrigins: []entities.Origin{},
		},
		"error on duplicate origin classes (class must be unique across all services)": {
			inputs: []input{
				{
					serviceID: "service1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
				{
					serviceID: "service2",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
			},
			wantErr: true,
		},
		"error on empty serviceID": {
			inputs: []input{
				{
					serviceID: "",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := pgtesting.NewDB(t)

			repo, err := NewOriginRepository(db)
			require.NoError(t, err)

			ctx := context.Background()

			var errs []error
			for _, input := range tt.inputs {
				err := repo.UpsertOrigins(ctx, input.serviceID, input.origins)
				errs = append(errs, err)
			}
			err = errors.Join(errs...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			// if all operarions were successful, verify final state
			require.NoError(t, err)
			gotOrigins, err := repo.ListOrigins(ctx)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.wantOrigins, gotOrigins) // order so far not guaranteed or relevant
		})
	}
}

func Test_UpsertOrigins_Concurrency(t *testing.T) {
	db := pgtesting.NewDB(t)

	repo, err := NewOriginRepository(db)
	require.NoError(t, err)

	serviceID := "some-service"
	origins := []entities.Origin{
		{Name: "origin4", Class: "classD", ServiceID: "read-only,ignored"},
		{Name: "origin3", Class: "classC", ServiceID: "read-only,ignored"},
	}

	var wg sync.WaitGroup
	iterations := 200

	for i := range iterations {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			err := repo.UpsertOrigins(context.Background(), serviceID, origins)
			if err != nil {
				t.Logf("Failed at iteration %d: %v", val, err)
			}
		}(i)
	}
	wg.Wait()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notification_service.origins WHERE service_id = $1", serviceID).Scan(&count)
	require.NoError(t, err)

	require.Equal(t, 2, count, "Data was duplicated due to race condition!")
}
