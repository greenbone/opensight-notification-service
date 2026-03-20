// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package notificationchannelservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/greenbone/opensight-notification-service/pkg/policy"
)

type TeamsService struct {
	transport *http.Client
}

func NewTeamsService(transport *http.Client) *TeamsService {
	return &TeamsService{transport: transport}
}

// SendMessage sends a message to the given MS Teams webhook URL.
// The message has to be in markdown format. For details see:
// https://learn.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features#markdown-commonmark-subset
func (m *TeamsService) SendMessage(webhookUrl string, message string) error {
	return sendTeamsMessage(m.transport, webhookUrl, message)
}

func sendTeamsMessage(transport *http.Client, webhookUrl, msg string) error {
	isTeamsOldWebhookUrl, err := policy.IsTeamsOldWebhookUrl(webhookUrl)
	if err != nil {
		return fmt.Errorf("failed to validate teams webhook url: %w", err)
	}

	var message map[string]interface{}
	if isTeamsOldWebhookUrl {
		message = map[string]interface{}{
			"text": msg,
		}

	} else {
		message = map[string]interface{}{
			"type": "message",
			"attachments": []map[string]interface{}{
				{
					"contentType": "application/vnd.microsoft.card.adaptive",
					"content": map[string]interface{}{
						"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
						"type":    "AdaptiveCard",
						"version": "1.2",
						"body": []map[string]interface{}{
							{
								"type": "TextBlock",
								"text": msg,
								"wrap": true,
							},
						},
					},
				},
			},
		}

	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("can not marshal teams message: %w", err)
	}

	resp, err := transport.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: timeout", ErrTeamsMessageDelivery)
		}
		return fmt.Errorf("%w: %w", ErrTeamsMessageDelivery, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: http status: %s", ErrTeamsMessageDelivery, resp.Status)
	}

	return nil
}
