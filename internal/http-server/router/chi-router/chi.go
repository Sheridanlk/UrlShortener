package chirouter

import (
	"UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/http-server/handlers/auth/login"
	"UrlShortener/internal/http-server/handlers/auth/register"
	"UrlShortener/internal/http-server/handlers/redirect"
	"UrlShortener/internal/http-server/handlers/url/save"
	"UrlShortener/internal/http-server/middleware/JWTAuth"
	"UrlShortener/internal/http-server/middleware/logger"
	"UrlShortener/internal/storage/postgresql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Setup(log *slog.Logger, storage *postgresql.Storage, ssoClient *grpc.Client, secret string) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/auth/register", register.New(log, ssoClient))
	router.Post("/auth/login", login.New(log, ssoClient))

	router.Route("/url", func(r chi.Router) {
		r.Use(JWTAuth.New(log, secret))

		r.Post("/create", save.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))

	return router
}
