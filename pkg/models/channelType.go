package models

type ChannelType string

const (
	ChannelTypeMail       ChannelType = "mail"
	ChannelTypeMattermost ChannelType = "mattermost"
	ChannelTypeTeams      ChannelType = "teams"
)

var AllowedChannels = []ChannelType{ChannelTypeMail, ChannelTypeMattermost, ChannelTypeTeams}

// HasRecipient returns true if the channel type requires/supports an explicit recipient.
func (ct ChannelType) HasRecipient() bool {
	return ct == ChannelTypeMail
}
