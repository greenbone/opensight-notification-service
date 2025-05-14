// SPDX-FileCopyrightText: 2025 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package models

type Sink struct {
	ID         string      `json:"id" readonly:"true"`
	Name       string      `json:"name" binding:"required"`
	Type       string      `json:"type" binding:"required" enums:"smtp, mattermost,teams"` // only populate `name` `type` and the root level poperty matching the type
	Mattermost *Mattermost `json:"mattermost"`
	MSTeams    *MSTeams    `json:"msteams"`
	SMTP       *SMTP       `json:"smtp"`
}

type Mattermost struct {
	Webhook string `json:"webhook" binding:"required"`
}

type MSTeams struct {
	Webhook string `json:"webhook" binding:"required"`
}

type SMTP struct {
	Host                 string `json:"host" binding:"required"`
	Port                 uint32 `json:"port" binding:"required"`
	UserName             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required"`
	Sender               string `json:"sender" binding:"required"`
	ConnectionSecurity   string `json:"connection_security" binding:"required" enums:"STARTTLS,SSL,NONE"`
	AuthenticationMethod string `json:"authentication_method" binding:"required" enums:"Plain,Encrypted,GSSAPI,Kerberos"`
}

type SinkReference struct {
	ID           string `json:"id" binding:"required"`
	Name         string `json:"name" readonly:"true"`
	Type         string `json:"type" readonly:"true"`
	HasRecipient bool   `json:"hasRecipient"` // indicates if the sink supports/requires specifying a specific recipient
}
