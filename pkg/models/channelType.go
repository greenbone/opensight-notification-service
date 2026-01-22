package models

type ChannelType string

const (
	ChannelTypeMail       ChannelType = "mail"
	ChannelTypeMattermost ChannelType = "mattermost"
	ChannelTypeTeams      ChannelType = "teams"
)
