package logging

import (
	"context"
	"fmt"
	"log/slog"
)

func String(key, val string) slog.Attr {
	return slog.String(key, val)
}

func Error(err error) slog.Attr {
	return String("error", err.Error())
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
		attributes = append(attributes, String("user_id", fmt.Sprint(userID)))
	}
	return attributes
}
