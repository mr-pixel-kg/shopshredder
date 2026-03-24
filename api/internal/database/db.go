package database

import (
	"log/slog"
	"time"

	"github.com/manuel/shopware-testenv-platform/api/internal/config"
	"github.com/manuel/shopware-testenv-platform/api/internal/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig, logLevel slog.Level) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logging.NewGormLogger(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
