package logging

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	"go.opentelemetry.io/otel/trace"
)

func String(key, val string) slog.Attr {
	return slog.String(key, val)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Error(err error) slog.Attr {
	var infrastructureError *errs.Error

	if errors.As(err, &infrastructureError) {
		attributes := []any{
			Int("code", int(infrastructureError.Code)),
			String("message", infrastructureError.Message),
		}
		if infrastructureError.Cause != nil {
			attributes = append(attributes, String("cause", infrastructureError.Cause.Error()))
		}
		return slog.Group("error", attributes...)
	}

	return slog.String("error", err.Error())
}

func traceAttributes(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	if !spanContext.IsValid() {
		return nil
	}

	return []slog.Attr{
		slog.String("trace_id", spanContext.TraceID().String()),
		slog.String("span_id", spanContext.SpanID().String()),
	}
}
