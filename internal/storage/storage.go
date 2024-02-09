package storage

import (
	"bot/internal/config"
	"bot/internal/storage/postgres"
	"bot/internal/storage/redis"
	"bot/pkg/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage struct {
	Redis    *redis.Redis
	Postgres *postgres.Postgres
}

func NewStorage(
	cfg *config.Config,
	redisClient *redis.Redis,
	postgresClient *postgres.Postgres,
) *Storage {
	setupStoragePool(&cfg.DB.Pool, postgresClient.Log, postgresClient.Db)

	return &Storage{
		Redis:    redisClient,
		Postgres: postgresClient,
	}
}

func setupStoragePool(cfg *config.DBPoolConfig, log *logging.Logger, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(cfg.IdleTimeout)
}
