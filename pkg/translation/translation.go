package translation

const (
	PortIsRequired            = "A port is required."
	UsernameIsRequired        = "A username is required."
	PasswordIsRequired        = "A password is required."
	ChannelNameIsRequired     = "A channel name is required."
	WebhookUrlIsRequired      = "A Webhook URL is required."
	ValidWebhookUrlIsRequired = "Please enter a valid webhook URL."

	// Email
	MailhubIsRequired          = "A mailhub is required."
	MailSenderIsRequired       = "A sender email is required."
	ValidEmailSenderIsRequired = "A valid sender email is required."

	// Mattermost
	MattermostChannelLimitReached     = "Mattermost channel limit reached."
	MattermostChannelNameAlreadyExist = "Mattermost channel name already exists."

	// Microsoft Teams
	TeamsChannelLimitReached     = "MS Teams channel limit reached."
	TeamsChannelNameAlreadyExist = "MS Teams channel name already exists."
)

// Rules
const (
	NameIsRequired        = "A name is required."
	LevelIsRequired       = "A level is required."
	OriginClassIsRequired = "An origin class is required."
	ChannelIsRequired     = "A channel is required."
	LevelsAreRequired     = "At least one level is required."
	OriginsAreRequired    = "At least one origin is required."

	RecipientRequiredForChannel     = "Recipient is required for the selected channel."
	RecipientNotSupportedForChannel = "Recipient is not supported for the selected channel."

	InvalidID        = "ID must be a valid UUIDv4."
	InvalidChannelID = "Channel ID must be a valid UUIDv4."

	RuleLimitReached      = "Alert rule limit reached."
	RuleNameAlreadyExists = "Alert rule name already exists."
	OriginsNotFound       = "One or more origins do not exist."
	ChannelNotFound       = "Channel does not exist."
)
