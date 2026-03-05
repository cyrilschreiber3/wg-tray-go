package logger

import (
	"log/slog"
	"os"
	"strings"
)

func InitLogger() *slog.Logger {
	logPath := getLogPath()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logFile = os.Stdout
	}
	return slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slogLevelFromEnv()}))
}

func slogLevelFromEnv() slog.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOGLEVEL"))) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
