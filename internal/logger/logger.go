package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var base *slog.Logger

func Init() {
	_ = godotenv.Load()

	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "unknown"
	}

	serviceName := os.Getenv("APP_NAME")
	if serviceName == "" {
		serviceName = "api-template"
	}

	var out io.Writer = os.Stdout
	if lokiURL := strings.TrimRight(os.Getenv("LOKI_URL"), "/"); lokiURL != "" {
		out = &lokiTee{url: lokiURL, serviceName: serviceName}
	}

	base = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level:       parseLevel(os.Getenv("LOG_LEVEL")),
		ReplaceAttr: replaceAttr,
	})).With(
		slog.String("service", serviceName),
		slog.String("hostname", host),
	)
}

func L() *slog.Logger {
	if base == nil {
		Init()
	}
	return base
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey && len(groups) == 0 {
		a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
		return a
	}
	switch strings.ToLower(a.Key) {
	case "password", "card", "token", "authorization", "secret":
		a.Value = slog.StringValue("[Redacted]")
	}
	return a
}

