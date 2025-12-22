package logger

import (
	"log/slog"
	"os"

	"github.com/hardal7/pex/internal/config"
	"github.com/lmittmann/tint"
)

func Init() {
	w := os.Stderr
	var level slog.Level
	switch config.LoggingType {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "SILENT":
		level = slog.LevelWarn
	default:
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: "15:04:05",
		}),
	))
}

func Debug(diagnostics string, args ...any) {
	slog.Debug(diagnostics, args...)
}

func Info(diagnostics string, args ...any) {
	slog.Info(diagnostics, args...)
}

func Warn(diagnostics string, args ...any) {
	slog.Warn(diagnostics, args...)
}

func Error(diagnostics string, args ...any) {
	slog.Error(diagnostics, args...)
	os.Exit(1)
}
