package JWTAuth

import (
	"UrlShortener/internal/config"
	"UrlShortener/internal/logger/sl"
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func New(log *slog.Logger, cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("component", "middleware/auth"),
			)

			ah := r.Header.Get("Authorization")
			if !strings.HasPrefix(ah, "Bearer ") {
				log.Info("invalid authentication method")

				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			tokenStr := strings.TrimPrefix(ah, "Bearer ")
			tokenParsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.AppSecret), nil
			})
			if err != nil {
				log.Error("unable to parse the token", sl.Err(err))

				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			claims, ok := tokenParsed.Claims.(jwt.MapClaims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			uid := int64(claims["uid"].(float64))
			ctx := context.WithValue(r.Context(), "uID", uid)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
