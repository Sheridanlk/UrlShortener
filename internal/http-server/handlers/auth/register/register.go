package register

import (
	"UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/lib/api/response"
	"UrlShortener/internal/logger/sl"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type Response struct {
	response.Response
	UserID int64 `json:"user_id"`
}

type RegisterService interface {
	Register(ctx context.Context, email string, password string) (int64, error)
}

func New(log *slog.Logger, registerService RegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.register"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded")

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		user_id, err := registerService.Register(context.Background(), req.Email, req.Password)

		if err != nil {
			if errors.Is(err, grpc.ErrUserExists) {
				log.Info("user is alredy exists", slog.String("email", req.Email))

				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, response.Error("user is already exists"))

				return
			}

			log.Info("failed to registrate user", slog.String("email", req.Email))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to registrate"))

			return
		}

		log.Info("user registrated", slog.Int64("user id", user_id))
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{
			UserID: user_id,
		})
	}
}
