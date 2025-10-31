package core

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yakuter/ugin/internal/config"
	"github.com/yakuter/ugin/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config, appLogger *logger.Logger) (*gorm.DB, error) {
	// Setup GORM logger
	dbLogLevel := gormlogger.Silent
	if cfg.Database.LogMode {
		dbLogLevel = gormlogger.Info
	}

	dbLogFile, err := os.OpenFile("ugin.db.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open database log file: %w", err)
	}

	gormLogger := gormlogger.New(
		log.New(dbLogFile, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  dbLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection based on driver
	var db *gorm.DB
	switch cfg.Database.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{
			Logger: gormLogger,
		})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
			Logger: gormLogger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
			Logger: gormLogger,
		})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	appLogger.Info("database connected", "driver", cfg.Database.Driver)
	return db, nil
}
