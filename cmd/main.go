package main

import (
	"UrlShortener/internal/config"
	"UrlShortener/internal/logger"
	"log/slog"
)

func main() {
	cfg := config.Load()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
}
