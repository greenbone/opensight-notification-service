package policy

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

func TeamsWebhookUrlPolicy(webhook string) (*url.URL, error) {
	if webhook == "" {
		return nil, errors.New("webhook URL is required")
	}

	u, err := url.Parse(webhook)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	var re = regexp.MustCompile(`^https://[\w.-]+/webhook/[a-zA-Z0-9]+$`)
	if !re.MatchString(webhook) {
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

	var re = regexp.MustCompile(`^https://[\w.-]+/hooks/[a-zA-Z0-9]+$`)
	if !re.MatchString(webhook) {
		return nil, errors.New("invalid Mattermost webhook URL")
	}

	return u, nil
}
