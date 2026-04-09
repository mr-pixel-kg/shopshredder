package database

import (
	"log/slog"
	"time"

	"github.com/mr-pixel-kg/shopshredder/api/internal/config"
	"github.com/mr-pixel-kg/shopshredder/api/internal/logging"
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
