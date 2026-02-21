package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Span(ctx context.Context, serviceName, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	return tracer.Start(ctx, spanName)
}
