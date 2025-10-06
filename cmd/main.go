package main

import (
	ssogrpc "UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/config"
	chirouter "UrlShortener/internal/http-server/router/chi-router"
	"UrlShortener/internal/http-server/server"
	"UrlShortener/internal/logger/sl"
	"UrlShortener/internal/storage/postgresql"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	log := sl.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Addres,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", sl.Err(err))
		os.Exit(1)
	}

	storage, err := postgresql.Init(cfg.PosgreSQL.User, cfg.PosgreSQL.Password, cfg.PosgreSQL.Name)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chirouter.Setup(log, storage, ssoClient, cfg.AppSecret)

	srv := server.New(log, router, cfg.HTTPServer.Address, cfg.HTTPServer.Timeout, cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout)

	go srv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	srv.Stop()
	storage.Close()

	log.Info("url-shortener stopped")
}
