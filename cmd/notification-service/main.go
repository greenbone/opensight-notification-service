// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller"

	"github.com/go-playground/validator"
	"github.com/greenbone/opensight-golang-libraries/pkg/logs"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/config/secretfiles"
	"github.com/greenbone/opensight-notification-service/pkg/repository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/healthservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/healthcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/notificationcontroller"
	"github.com/kelseyhightower/envconfig"

	"github.com/rs/zerolog/log"
)

func main() {
	var cfg config.Config
	// Note: secrets can be passed directly by env var or via file
	// if the same secret is supplied in both ways, the env var takes precedence
	err := secretfiles.Read(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read secrets from files")
	}
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}
	err = validator.New().Struct(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}
	err = logs.SetupLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up logger")
	}

	check(run(cfg))
}

func run(config config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	pgClient, err := repository.NewClient(config.Database)
	if err != nil {
		return err
	}

	// auth
	realmInfo := auth.KeycloakRealmInfo{
		RealmId:               config.KeycloakConfig.Realm,
		AuthServerInternalUrl: config.KeycloakConfig.AuthServerUrl,
	}

	authorizer, err := auth.NewKeycloakAuthorizer(realmInfo)
	if err != nil {
		return fmt.Errorf("error creating keycloak token authorizer: %w", err)
	}

	authMiddleware, err := auth.NewGinAuthMiddleware(authorizer.ParseRequest)
	if err != nil {
		return fmt.Errorf("error creating keycloak auth middleware: %w", err)
	}

	notificationRepository, err := notificationrepository.NewNotificationRepository(pgClient)
	if err != nil {
		return err
	}

	notificationChannelRepository, err := notificationrepository.NewNotificationChannelRepository(pgClient)
	if err != nil {
		return err
	}

	notificationService := notificationservice.NewNotificationService(notificationRepository)
	notificationChannelService := notificationchannelservice.NewNotificationChannelService(notificationChannelRepository)
	healthService := healthservice.NewHealthService(pgClient)

	gin := web.NewWebEngine(config.Http)
	rootRouter := gin.Group("/")
	notificationServiceRouter := gin.Group("/api/notification-service")
	docsRouter := gin.Group("/docs/notification-service")

	// rest api docs
	web.RegisterSwaggerDocsRoute(docsRouter, config.KeycloakConfig)
	healthcontroller.RegisterSwaggerDocsRoute(docsRouter, config.KeycloakConfig)

	//instantiate controllers
	notificationcontroller.NewNotificationController(notificationServiceRouter, notificationService, authMiddleware)
	mailcontroller.NewMailController(notificationServiceRouter, notificationChannelService, authMiddleware)
	healthcontroller.NewHealthController(rootRouter, healthService) // for health probes (not a data source)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Http.Port),
		Handler:      gin,
		ReadTimeout:  config.Http.ReadTimeout,
		WriteTimeout: config.Http.WriteTimeout,
		IdleTimeout:  config.Http.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			check(err)
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Received signal. Shutting down")
	err = srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shut down error: %w", err)
	}

	return nil
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("critical error")
	}
}
