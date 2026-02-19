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

func (l *Logger) Debug(ctx context.Context, msg string, attributes ...slog.Attr) {
	attributes = append(attributes, traceAttributes(ctx)...)
	l.logger.LogAttrs(ctx, slog.LevelDebug, msg, attributes...)
}

func (l *Logger) Info(ctx context.Context, msg string, attributes ...slog.Attr) {
	attributes = append(attributes, traceAttributes(ctx)...)
	l.logger.LogAttrs(ctx, slog.LevelInfo, msg, attributes...)
}

func (l *Logger) Warn(ctx context.Context, msg string, attributes ...slog.Attr) {
	attributes = append(attributes, traceAttributes(ctx)...)
	l.logger.LogAttrs(ctx, slog.LevelWarn, msg, attributes...)
}

func (l *Logger) Error(ctx context.Context, msg string, attributes ...slog.Attr) {
	attributes = append(attributes, traceAttributes(ctx)...)
	l.logger.LogAttrs(ctx, slog.LevelError, msg, attributes...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, attributes ...slog.Attr) {
	attributes = append(attributes, traceAttributes(ctx)...)
	l.logger.LogAttrs(ctx, slog.LevelError, msg, attributes...)
	os.Exit(1)
}
