package asynqlog

import (
	"fmt"
	"log/slog"
)

type Logger struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Logger {
	return &Logger{log: log.With(slog.String("component", "asynq"))}
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(fmt.Sprint(args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(fmt.Sprint(args...))
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Error(fmt.Sprint(args...))
}
