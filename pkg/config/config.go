// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/configReader"
)

type Config struct {
	Http     Http
	Database Database
	LogLevel string `viperEnv:"LOG_LEVEL" default:"info"`
}

type Http struct {
	Port         int           `validate:"required,min=1,max=65535" viperEnv:"HTTP_PORT" default:"8085"`
	ReadTimeout  time.Duration `viperEnv:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `viperEnv:"HTTP_WRITE_TIMEOUT" default:"60s"`
	IdleTimeout  time.Duration `viperEnv:"HTTP_IDLE_TIMEOUT" default:"60s"`
}

type Database struct {
	Host     string `viperEnv:"DB_HOST" default:"localhost"`
	Port     int    `validate:"required,min=1,max=65535" viperEnv:"DB_PORT" default:"5432"`
	User     string `validate:"required" viperEnv:"DB_USERNAME"`
	Password string `validate:"required" viperEnv:"DB_PASSWORD"`
	DBName   string `validate:"required" viperEnv:"DB_NAME"`
	SSLMode  string `viperEnv:"DB_SSL_MODE" default:"require"`
}

func ReadConfig() (Config, error) {
	var config Config
	_, err := configReader.ReadEnvVarsIntoStruct(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
