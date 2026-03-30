// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

import (
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/notifications"
	"github.com/stretchr/testify/require"
)

func Test_RuleIsTriggered(t *testing.T) {

	// notification processed in all tests
	notification := Notification{
		Origin:      "Test Origin",
		OriginClass: "/serviceID/origin1",
		Timestamp:   "2024-01-01T00:00:00Z",
		Title:       "Test Notification",
		Detail:      "This is a test notification",
		Level:       notifications.LevelInfo,
	}

	tests := map[string]struct {
		rule Rule
		want bool
	}{
		"matching origin class and level triggers": {
			rule: ruleValid(func(r *Rule) {
				r.Trigger = Trigger{
					Origins: []OriginReference{{Class: notification.OriginClass}},
					Levels:  []notifications.Level{notification.Level},
				}
			}),
			want: true,
		},
		"deactivated rule does not trigger": {
			rule: ruleValid(func(r *Rule) {
				r.Active = false
			}),
			want: false,
		},
		"origin 'Any' always matches": {
			rule: ruleValid(func(r *Rule) {
				r.Trigger = Trigger{
					Origins: []OriginReference{{Class: OriginAllClass}},
					Levels:  []notifications.Level{notification.Level},
				}
			}),
			want: true,
		},
		"non-matching origin class does not trigger": {
			rule: ruleValid(func(r *Rule) {
				r.Trigger = Trigger{
					Origins: []OriginReference{{Class: "no-match"}},
					Levels:  []notifications.Level{notification.Level},
				}
			}),
			want: false,
		},
		"non-matching level does not trigger": {
			rule: ruleValid(func(r *Rule) {
				r.Trigger = Trigger{
					Origins: []OriginReference{{Class: notification.OriginClass}},
					Levels:  []notifications.Level{notifications.LevelError},
				}
			}),
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.rule.IsTriggered(notification)
			require.Equal(t, tt.want, got)
		})
	}
}

func ruleValid(options ...func(*Rule)) Rule {
	rule := Rule{
		Name: "Test Rule",
		Trigger: Trigger{
			Origins: []OriginReference{{Class: "no-match"}},
			Levels:  []notifications.Level{notifications.LevelInfo},
		},
		Action: Action{
			Channel: ChannelReference{
				ID:   "00000000-0000-0000-0000-000000000000",
				Type: ChannelTypeMattermost,
			},
		},
		Active: true,
	}

	for _, option := range options {
		if option != nil {
			option(&rule)
		}
	}

	return rule
}
