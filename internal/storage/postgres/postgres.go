package postgres

import (
	"bot/internal/config"
	"bot/internal/storage/postgres/repositories"
	"bot/pkg/logging"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	Log  *logging.Logger
	Db   *gorm.DB
	User *repositories.UserRepository
}

func NewPostgres(cfg *config.DBConfig, log *logging.Logger) *Postgres {
	log.Info("Connecting to database...")
	dbDialectName := dbDialect(cfg)
	db := dbConnect(log, dbDialectName)
	db.Logger = logger.Default.LogMode(logger.Silent)

	return &Postgres{
		Log:  log,
		Db:   db,
		User: repositories.NewUserRepository(log, db),
	}
}

func dbConnect(log *logging.Logger, dbDialect gorm.Dialector) *gorm.DB {
	dbConn, err := gorm.Open(dbDialect, &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	return dbConn
}

func dbDialect(cfg *config.DBConfig) gorm.Dialector {
	dsn := cfg.ConnString()
	return postgres.Open(dsn)
}
