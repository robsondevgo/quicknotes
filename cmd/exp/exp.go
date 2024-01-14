package main

import (
	"log/slog"
	"os"
)

func main() {
	// h := slog.NewTextHandler(os.Stderr, nil)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	log := slog.New(h).With("app", "exp")

	log.Debug("debug message")
	log.Info("info message", slog.Group("request_info", "request_id", 1, "user", "robson"))
	log.Warn("warn message")
	log.Error("error message", "request_id", 1)
}
