package main

import (
	"context"
	"errors"
	"flag"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appconfig "kbfood/internal/config"
	"kbfood/internal/domain/service"
	dbinfra "kbfood/internal/infra/db"
	"kbfood/internal/infra/platform"
	repoimpl "kbfood/internal/infra/repository"
	schedulerinfra "kbfood/internal/infra/scheduler"
	httpiface "kbfood/internal/interface/http"
	"kbfood/internal/interface/http/handler"
	applog "kbfood/internal/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg, err := appconfig.Load(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := applog.Init(&applog.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     "stdout",
		TimeFormat: cfg.Log.TimeFormat,
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize logger")
	}

	ctx := context.Background()
	database, err := dbinfra.NewPool(ctx, &dbinfra.Config{
		Path:            cfg.Database.Path,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer database.Close()

	queries := database.Queries()

	productRepo := repoimpl.NewProductRepository(queries)
	masterProductRepo := repoimpl.NewMasterProductRepository(queries)
	notificationRepo := repoimpl.NewNotificationRepository(queries)
	blockedRepo := repoimpl.NewBlockedRepository(queries)
	trendRepo := repoimpl.NewTrendRepository(queries)
	candidateRepo := repoimpl.NewCandidateRepository(queries)
	userSettingsRepo := repoimpl.NewUserSettingsRepository(queries)
	syncStatusRepo := repoimpl.NewSyncStatusRepository(database)

	cleaningService := service.NewDataCleaningService(masterProductRepo, candidateRepo, trendRepo)
	notificationService := service.NewNotificationService(
		notificationRepo,
		productRepo,
		masterProductRepo,
		userSettingsRepo,
		cfg.BarkURL,
	)

	tttClient := platform.NewTanTanTangClient(&cfg.Platforms.TanTanTang)

	syncJob := schedulerinfra.NewSyncJob(cfg, tttClient, cleaningService, syncStatusRepo)
	promoteCandidatesJob := schedulerinfra.NewPromoteCandidatesJob(cleaningService)
	priceCheckJob := schedulerinfra.NewPriceCheckJob(notificationService)
	recordTrendsJob := schedulerinfra.NewRecordTrendsJob(cleaningService)

	scheduler := schedulerinfra.NewScheduler(nil)
	registerJob(scheduler, syncJob, "0 */5 * * * *")
	registerJob(scheduler, priceCheckJob, "0 */5 * * * *")
	registerJob(scheduler, promoteCandidatesJob, "*/30 * * * * *")
	registerJob(scheduler, recordTrendsJob, "0 5 0 * * *")
	scheduler.Start()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := scheduler.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("failed to stop scheduler cleanly")
		}
	}()

	productHandler := handler.NewProductHandler(productRepo, masterProductRepo, notificationRepo, blockedRepo, trendRepo)
	externalHandler := handler.NewExternalHandler(cleaningService)
	syncHandler := handler.NewSyncHandler(syncJob, tttClient)
	statusHandler := handler.NewStatusHandler(syncStatusRepo)
	userHandler := handler.NewUserHandler(userSettingsRepo)

	router := httpiface.Router(
		productHandler,
		externalHandler,
		syncHandler,
		statusHandler,
		userHandler,
		database,
	)

	server := &stdhttp.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.ReadTimeout,
	}

	go func() {
		log.Info().
			Str("addr", cfg.Server.Address()).
			Str("dbPath", cfg.Database.Path).
			Msg("starting server")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server stopped unexpectedly")
		}
	}()

	waitForShutdown(server, cfg.Server.ShutdownTimeout)
}

func registerJob(scheduler *schedulerinfra.Scheduler, job schedulerinfra.Job, cronExpr string) {
	if err := scheduler.RegisterJob(job, cronExpr); err != nil {
		log.Fatal().
			Err(err).
			Str("job", job.Name()).
			Str("cron", cronExpr).
			Msg("failed to register job")
	}
}

func waitForShutdown(server *stdhttp.Server, timeout time.Duration) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shut down server cleanly")
		return
	}

	log.Info().Msg("server shut down cleanly")
}
