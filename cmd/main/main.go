package main

import (
	"bot/internal/bot"
	"bot/internal/config"
	"bot/internal/storage"
	myPostgres "bot/internal/storage/postgres"
	myRedis "bot/internal/storage/redis"
	"bot/pkg/logging"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const (
	levelsUpToRootDir = 3
)

func main() {
	cfg := config.MustLoad(levelsUpToRootDir)

	log := setupLogger(
		cfg.Env,
		config.StructDateFormat,
		cfg.Paths.ConfigInfoPath,
		cfg.Paths.ConfigDebugPath,
	)

	storageClient := storage.NewStorage(
		cfg,
		myRedis.NewRedis(&cfg.RedisDB, log),
		myPostgres.NewPostgres(&cfg.DB, log),
	)

	go bot.NewBot(cfg, log, storageClient).Run()

	// Gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("Got signal", zap.String("signal", sign.String()))
	log.Info("Shutting down...")
}

func setupLogger(env, structDateFormat, pathToInfoLogs, pathToDebugLogs string) *logging.Logger {
	log := logging.NewLogger(
		logging.InitLogger(
			env,
			structDateFormat,
			pathToInfoLogs,
			pathToDebugLogs,
		),
	)

	log.Info("Initializing logger and config...")
	log.Debug("Debug mode on...")

	go logging.ClearLogFiles(
		pathToInfoLogs,
		pathToDebugLogs,
		structDateFormat,
		log,
	)

	return log
}
