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
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/greenbone/keycloak-client-golang/auth"
	"github.com/greenbone/opensight-notification-service/pkg/security"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationchannelservice"
	"github.com/greenbone/opensight-notification-service/pkg/web/errmap"
	"github.com/greenbone/opensight-notification-service/pkg/web/mailcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/teamsController"
	"github.com/jmoiron/sqlx"

	"github.com/go-playground/validator"
	"github.com/greenbone/opensight-notification-service/pkg/jobs/checkmailconnectivity"
	"github.com/greenbone/opensight-notification-service/pkg/web/mattermostcontroller"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

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
	defer func(pgClient *sqlx.DB) {
		err = pgClient.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Error while closing Postgres: %s", err)
		}
	}(pgClient)
	log.Debug().Msg("postgres database connection successful")

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
		return fmt.Errorf("error creating Notification Repository: %w", err)
	}

	// Encrypt
	manager := security.NewEncryptManager()
	manager.UpdateKeys(config.DatabaseEncryptionKey)

	notificationChannelRepository, err := notificationrepository.NewNotificationChannelRepository(pgClient, manager)
	if err != nil {
		return fmt.Errorf("error creating Notification Channel Repository: %w", err)
	}

	mailService := notificationchannelservice.NewMailService()
	notificationService := notificationservice.NewNotificationService(notificationRepository)
	notificationChannelService := notificationchannelservice.NewNotificationChannelService(notificationChannelRepository)
	mailChannelService := notificationchannelservice.NewMailChannelService(notificationChannelService, notificationChannelRepository, mailService, config.ChannelLimit.EMailLimit)
	mattermostChannelService := notificationchannelservice.NewMattermostChannelService(notificationChannelService, config.ChannelLimit.MattermostLimit)

	notificationTransport := http.Client{Timeout: 15 * time.Second}
	teamsChannelService := notificationchannelservice.NewTeamsChannelService(
		notificationChannelService, config.ChannelLimit.TeamsLimit, notificationTransport)
	healthService := healthservice.NewHealthService(pgClient)

	// scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("error creating scheduler: %w", err)
	}
	_, err = scheduler.NewJob(
		gocron.DurationJob(1*time.Hour),
		gocron.NewTask(checkmailconnectivity.NewJob(notificationService, notificationChannelService, mailChannelService)),
	)
	if err != nil {
		return fmt.Errorf("error creating mail connectivity check job: %w", err)
	}
	scheduler.Start()

	registry := errmap.NewRegistry()

	router := web.NewWebEngine(config.Http, registry)

	// rest api docs
	docsRouter := router.Group("/docs/notification-service")
	web.RegisterSwaggerDocsRoute(docsRouter, config.KeycloakConfig)
	healthcontroller.RegisterSwaggerDocsRoute(docsRouter, config.KeycloakConfig)

	// instantiate controllers
	notificationServiceRouter := router.Group("/api/notification-service")
	notificationcontroller.AddNotificationController(notificationServiceRouter, notificationService, authMiddleware)
	mailcontroller.NewMailController(notificationServiceRouter, notificationChannelService, mailChannelService, authMiddleware, registry)
	mailcontroller.AddCheckMailServerController(notificationServiceRouter, mailChannelService, authMiddleware, registry)
	mattermostcontroller.NewMattermostController(notificationServiceRouter, notificationChannelService, mattermostChannelService, authMiddleware, registry)
	teamsController.AddTeamsController(notificationServiceRouter, notificationChannelRepository, teamsChannelService, authMiddleware, registry)

	// health router
	rootRouter := router.Group("/")
	healthcontroller.NewHealthController(rootRouter, healthService) // for health probes (not a data source)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Http.Port),
		Handler:      router,
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
