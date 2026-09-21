package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Env string `default:"local" envconfig:"LOGGER_ENV"`
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(c Config) *slog.Logger {
	var log *slog.Logger

	switch c.Env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return log
}
