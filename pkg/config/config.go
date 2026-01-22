// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"time"
)

// Note: The `envconfig` entries of fields in nested structs are prefixed by the `envconfig` entry of the struct
// I.e. the env var in {INNER1:{INNER2:{FIELD1:"foo"}}} for FIELD1 is `INNER1_INNER2_FIELD1`

type Config struct {
	Http                  Http                  `envconfig:"HTTP"`
	Database              Database              `envconfig:"DB"`
	LogLevel              string                `envconfig:"LOG_LEVEL" default:"info"`
	KeycloakConfig        KeycloakConfig        `envconfig:"KEYCLOAK"`
	ChannelLimit          ChannelLimits         `envconfig:"CHANNELLIMIT"`
	DatabaseEncryptionKey DatabaseEncryptionKey `envconfig:"DATABASE_ENCRYPTION_KEY"`
}

type ChannelLimits struct {
	EMailLimit      int `envconfig:"EMAIL_LIMIT" default:"1"`
	MattermostLimit int `envconfig:"MATTERMOST_LIMIT" default:"20"`
}

type Http struct {
	Port           int           `validate:"required,min=1,max=65535" envconfig:"PORT" default:"8085"`
	ReadTimeout    time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout   time.Duration `envconfig:"WRITE_TIMEOUT" default:"60s"`
	IdleTimeout    time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
	AllowedOrigins []string      `envconfig:"ALLOWED_ORIGINS" default:"https://opensight-lookout.greenbone.io"`
}

type Database struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `validate:"required,min=1,max=65535" envconfig:"PORT" default:"5432"`
	User     string `validate:"required" envconfig:"USERNAME"`
	Password string `validate:"required" envconfig:"PASSWORD"`
	DBName   string `validate:"required" envconfig:"NAME"`
	SSLMode  string `envconfig:"SSL_MODE" default:"require"`
}

type DatabaseEncryptionKey struct {
	Password     string `envconfig:"DATABASE_ENCRYPTION_KEY_PASSWORD" default:"database-encryption-key-default-password"`
	PasswordSalt string `envconfig:"DATABASE_ENCRYPTION_KEY_PASSWORD_SALT" default:"database-encryption-key-default-password-salt"`
}

type KeycloakConfig struct {
	Realm         string `validate:"required" envconfig:"REALM" default:"opensight"`
	AuthServerUrl string `validate:"required" envconfig:"URL" default:"http://localhost:8082/auth"`
	WebClientName string `validate:"required" envconfig:"WEBCLIENT_NAME" default:"local-web"`
	PublicUrl     string `validate:"required" envconfig:"PUBLIC_URL" default:"http://localhost:8082/auth"`
}
