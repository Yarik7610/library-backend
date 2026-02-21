package logging

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func String(key, val string) slog.Attr {
	return slog.String(key, val)
}

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func Error(err error) slog.Attr {
	return String("error", err.Error())
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
