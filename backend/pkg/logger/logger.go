package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// Setup configures the default slog logger to output structured JSON to stdout.
// Sevalla parses each JSON line for its log viewer, populating severity from the
// "level" field and attributes from all other keys.
// Log records are enriched with trace_id and span_id when OTel context is present.
// Must be called early in main() before any slog calls.
func Setup() {
	json := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	slog.SetDefault(slog.New(&traceHandler{Handler: json}))
}

// traceHandler wraps an slog.Handler and adds OTel trace context fields
// (trace_id, span_id) to every log record that carries a valid span context.
type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name)}
}
