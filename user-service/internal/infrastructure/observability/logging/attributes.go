package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/errs"
)

func String(key, val string) slog.Attr {
	return slog.String(key, val)
}

func Any(key string, val any) slog.Attr {
	return slog.Any(key, val)
}

func Error(err error) slog.Attr {
	var infrastructureError *errs.Error

	if errors.As(err, &infrastructureError) {
		attributes := []any{
			slog.Int("code", int(infrastructureError.Code)),
			slog.String("message", infrastructureError.Message),
		}
		if infrastructureError.Cause != nil {
			attributes = append(attributes, slog.String("cause", infrastructureError.Cause.Error()))
		}
		return slog.Group("error", attributes...)
	}

	return slog.Group("error", slog.String("message", infrastructureError.Cause.Error()))
}

func TraceAttributes(ctx context.Context) []slog.Attr {
	var attributes []slog.Attr

	if traceID := ctx.Value("trace_id"); traceID != nil {
		attributes = append(attributes, String("trace_id", fmt.Sprint(traceID)))
	}
	if spanID := ctx.Value("span_id"); spanID != nil {
		attributes = append(attributes, String("span_id", fmt.Sprint(spanID)))
	}
	if userID := ctx.Value("user_id"); userID != nil {
	}
	return attributes
}
