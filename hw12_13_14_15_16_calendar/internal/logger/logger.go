package logger

import (
	"log/slog"
	"strings"
)

type Logger struct {
	level  string
	logger *slog.Logger
}

func New(level string) *Logger {
	normalizedLevel := strings.ToLower(level)

	switch normalizedLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	return &Logger{level, slog.Default()}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}
