package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	log *slog.Logger
}

func NewSlogLogger(filePath, level string) (Logger, error) {

	out, err := output(filePath)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource: true,
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

func output(filePath string) (io.Writer, error) {
	if filePath == "" {
		return os.Stdout, nil
	} else {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}
