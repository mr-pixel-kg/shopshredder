package logging

import (
	"log/slog"
	"os"
	"time"

	charmlog "github.com/charmbracelet/log"
	"github.com/manuel/shopware-testenv-platform/api/internal/config"
	"github.com/muesli/termenv"
)

var TextMode bool

func Setup(cfg config.LoggingConfig) {
	TextMode = cfg.Format == config.LogFormatText
	level := ParseLevel(cfg.Level)

	switch cfg.Format {
	case config.LogFormatText:
		charm := charmlog.NewWithOptions(os.Stderr, charmlog.Options{
			Level:           charmlog.Level(level),
			ReportCaller:    level <= slog.LevelDebug,
			ReportTimestamp: true,
			TimeFormat:      time.TimeOnly,
		})
		charm.SetColorProfile(termenv.ANSI256)
		charmlog.SetDefault(charm)
		slog.SetDefault(slog.New(charm))

	default:
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: level <= slog.LevelDebug,
		})))
	}
}

func ParseLevel(l config.LogLevel) slog.Level {
	switch l {
	case config.LogLevelDebug:
		return slog.LevelDebug
	case config.LogLevelWarn:
		return slog.LevelWarn
	case config.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
