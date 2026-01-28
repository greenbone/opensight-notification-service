package policy

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var teamsRegex = regexp.MustCompile(`^https://[\w.-]+/webhook/[a-zA-Z0-9]+$`)
var mattermostRegex = regexp.MustCompile(`^https://[\w.-]+/hooks/[a-zA-Z0-9]+$`)

func TeamsWebhookUrlPolicy(webhook string) (*url.URL, error) {
	if webhook == "" {
		return nil, errors.New("webhook URL is required")
	}

	u, err := url.Parse(webhook)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if !teamsRegex.MatchString(webhook) {
		return nil, errors.New("invalid Teams webhook URL")
	}

	return u, nil
}

func MattermostWebhookUrlPolicy(webhook string) (*url.URL, error) {
	if webhook == "" {
		return nil, errors.New("webhook URL is required")
	}

	u, err := url.Parse(webhook)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if !mattermostRegex.MatchString(webhook) {
		return nil, errors.New("invalid Mattermost webhook URL")
	}

	return u, nil
}
