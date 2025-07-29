package logger

import (
	"log"
	"log/slog"
	"strings"
)

type logger struct {
	log *slog.Logger
}

var _ Logger = (*logger)(nil)

func NewSlogLogger(filePath, level string) (*logger, error) {
	out, err := output(filePath)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource: false,
		Level:     slogLevel(level),
	})

	return &logger{
		log: slog.New(handler),
	}, nil
}

func (l *logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

func (l *logger) ServerErrorLog() *log.Logger {
	return slog.NewLogLogger(l.log.Handler(), slog.LevelError)
}

func slogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
