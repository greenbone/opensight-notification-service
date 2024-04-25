package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/greenbone/opensight-notification-service/pkg/logging"
	"github.com/greenbone/opensight-notification-service/pkg/repository"
	"github.com/greenbone/opensight-notification-service/pkg/repository/notificationrepository"
	"github.com/greenbone/opensight-notification-service/pkg/services/healthservice"
	"github.com/greenbone/opensight-notification-service/pkg/services/notificationservice"
	"github.com/greenbone/opensight-notification-service/pkg/web"
	"github.com/greenbone/opensight-notification-service/pkg/web/healthcontroller"
	"github.com/greenbone/opensight-notification-service/pkg/web/notificationcontroller"

	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	err = logging.SetupLogger(config.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up logger")
	}

	check(run(config))
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
	notificationRepository, err := notificationrepository.NewNotificationRepository(pgClient)
	if err != nil {
		return err
	}

	notificationService := notificationservice.NewNotificationService(notificationRepository)
	healthService := healthservice.NewHealthService(pgClient)

	gin := web.NewWebEngine()
	rootRouter := gin.Group("/")
	notificationServiceRouter := gin.Group("/api/notification-service")
	docsRouter := gin.Group("/docs/notification-service")

	// rest api docs
	web.RegisterSwaggerDocsRoute(docsRouter)
	healthcontroller.RegisterSwaggerDocsRoute(docsRouter)

	//instantiate controllers
	notificationcontroller.NewNotificationController(notificationServiceRouter, notificationService)
	healthcontroller.NewHealthController(rootRouter, healthService) // for health probes (not a data source)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Http.Port),
		Handler:      gin,
		ReadTimeout:  config.Http.ReadTimeout,
		WriteTimeout: config.Http.WriteTimeout,
		IdleTimeout:  config.Http.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
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
