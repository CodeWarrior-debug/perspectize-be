package logger

import (
	"log/slog"
	"os"
)

// Setup configures the default slog logger to output structured JSON to stdout.
// Sevalla parses each JSON line for its log viewer, populating severity from the
// "level" field and attributes from all other keys.
// Must be called early in main() before any slog calls.
func Setup() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler))
}
