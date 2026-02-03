// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import (
	"context"
	"testing"

	"github.com/greenbone/opensight-notification-service/pkg/entities"
	"github.com/greenbone/opensight-notification-service/pkg/pgtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateOrigins_ListOrigins(t *testing.T) {
	type input struct {
		origins   []entities.Origin
		namespace string
	}

	tests := map[string]struct {
		inputs      []input
		wantOrigins []entities.Origin
		wantErr     bool
	}{
		"create origins in single namespace (no prior data)": {
			inputs: []input{
				{
					namespace: "ns1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA", Namespace: "read only, to be ignored"},
						{Name: "origin2", Class: "classB"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin1", Class: "classA", Namespace: "ns1"},
				{Name: "origin2", Class: "classB", Namespace: "ns1"},
			},
		},
		"create origins in multiple namespaces": {
			inputs: []input{
				{
					namespace: "ns1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
				{
					namespace: "ns2",
					origins: []entities.Origin{
						{Name: "origin2", Class: "classB"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin1", Class: "classA", Namespace: "ns1"},
				{Name: "origin2", Class: "classB", Namespace: "ns2"},
			},
		},
		"upsert origins replaces the entries from the same namespace": {
			inputs: []input{
				{
					namespace: "ns1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
						{Name: "origin2", Class: "classB"},
					},
				},
				{
					namespace: "ns2",
					origins: []entities.Origin{
						{Name: "origin3", Class: "classC"},
					},
				},
				{
					namespace: "ns1",
					origins: []entities.Origin{
						{Name: "origin4", Class: "classD"},
					},
				},
			},
			wantOrigins: []entities.Origin{
				{Name: "origin4", Class: "classD", Namespace: "ns1"},
				{Name: "origin3", Class: "classC", Namespace: "ns2"},
			},
		},
		"delete all origins in namespace when upserting empty list": {
			inputs: []input{
				{
					namespace: "ns1",
					origins: []entities.Origin{
						{Name: "origin1", Class: "classA"},
					},
				},
				{
					namespace: "ns1",
					origins:   []entities.Origin{},
				},
			},
			wantOrigins: []entities.Origin{},
		},
		"error on empty namespace": {
			inputs: []input{
				{
					namespace: "",
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

			for _, input := range tt.inputs {
				err := repo.UpsertOrigins(ctx, input.namespace, input.origins)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}

			// if all operarions were successful, verify final state
			gotOrigins, err := repo.ListOrigins(ctx)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.wantOrigins, gotOrigins) // order so far not guaranteed or relevant
		})
	}
}
