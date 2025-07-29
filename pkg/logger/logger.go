package logger

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	ServerErrorLog() *log.Logger
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
