package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/vanshsharma3777/WorkerBox/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dsn string) error {
	var err error
	var sqlDB *sql.DB

	// Retry loop — useful for Neon cold starts
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err = DB.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					slog.Info(
						"connected to PostgreSQL",
						"attempt",
						i+1,
					)
					break
				}
			}
		}

		slog.Warn(
			"waiting for database...",
			"attempt",
			i+1,
			"error",
			err,
		)

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	// Connection pool configuration
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	slog.Info("running database migrations")

	err := DB.AutoMigrate(
		&models.Job{},
	)

	if err != nil {
		return err
	}

	slog.Info("database migrated successfully")

	return nil
}
