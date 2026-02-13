package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(env string) *Logger {
	env = strings.ToLower(env)

	var loggerHandler slog.Handler

	if env == "local" {
		loggerHandler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		})
	} else {
		loggerHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelWarn,
		})
	}

	return &Logger{logger: slog.New(loggerHandler)}
}

func (l *Logger) Debug(msg string, attributes ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attributes...)
}

func (l *Logger) Info(msg string, attributes ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attributes...)
}

func (l *Logger) Warn(msg string, attributes ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg, attributes...)
}

func (l *Logger) Error(msg string, attributes ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, attributes...)
}

func (l *Logger) Fatal(msg string, attributes ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, attributes...)
	os.Exit(1)
}

func (l *Logger) With(attributes ...slog.Attr) *Logger {
	updatedLoggerHandler := l.logger.Handler().WithAttrs(attributes)
	return &Logger{logger: slog.New(updatedLoggerHandler)}
}
