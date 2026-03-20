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
)

type MattermostService struct {
	transport *http.Client
}

func NewMattermostService(transport *http.Client) *MattermostService {
	return &MattermostService{transport: transport}
}

// SendMessage sends a message to the given Mattermost webhook URL.
// The message has to be in markdown format. For details see:
// https://docs.mattermost.com/end-user-guide/collaborate/format-messages.html#use-markdown
func (m *MattermostService) SendMessage(webhookUrl string, message string) error {
	return sendMattermostMessage(m.transport, webhookUrl, message)
}

func sendMattermostMessage(transport *http.Client, webhookUrl string, message string) error {
	body, err := json.Marshal(map[string]string{
		"text": message,
	})
	if err != nil {
		return fmt.Errorf("can not marshal mattermost message: %w", err)
	}

	resp, err := transport.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: timeout", ErrMattermostMassageDelivery)
		}
		return fmt.Errorf("%w: %w", ErrMattermostMassageDelivery, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: http status: %s", ErrMattermostMassageDelivery, resp.Status)
	}

	return nil
}
