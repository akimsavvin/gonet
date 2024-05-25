package logging

import (
	"context"
	"log/slog"
)

type middleware struct {
	next slog.Handler
}

func newMiddleware(next slog.Handler) *middleware {
	return &middleware{next: next}
}

func (m *middleware) Handle(ctx context.Context, rec slog.Record) error {
	if val, ok := ctx.Value("request_id").(string); ok && val != "" {
		rec.Add(slog.String("request_id", val))
	}

	if val, ok := ctx.Value("correlation_id").(string); ok && val != "" {
		rec.Add(slog.String("correlation_id", val))
	}

	return m.next.Handle(ctx, rec)
}

func (m *middleware) Enabled(ctx context.Context, level slog.Level) bool {
	return m.next.Enabled(ctx, level)
}

func (m *middleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newMiddleware(m.next.WithAttrs(attrs))
}

func (m *middleware) WithGroup(name string) slog.Handler {
	return newMiddleware(m.next.WithGroup(name))
}
