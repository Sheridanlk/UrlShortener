package main

import (
	ssogrpc "UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/config"
	"UrlShortener/internal/http-server/handlers/auth/login"
	"UrlShortener/internal/http-server/handlers/auth/register"
	"UrlShortener/internal/http-server/handlers/redirect"
	"UrlShortener/internal/http-server/handlers/url/save"
	"UrlShortener/internal/http-server/middleware/JWTAuth"
	"UrlShortener/internal/http-server/middleware/logger"
	"UrlShortener/internal/logger/sl"
	"UrlShortener/internal/storage/postgresql"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/auth/register", register.New(log, ssoClient))
	router.Post("/auth/login", login.New(log, ssoClient))

	router.Route("/url", func(r chi.Router) {
		r.Use(JWTAuth.New(log, *cfg))

		r.Post("/create", save.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("addres", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
