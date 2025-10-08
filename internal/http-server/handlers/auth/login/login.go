package login

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
	Password string `json:"password" validate:"required"`
	AppID    int32  `json:"app_id" validate:"required"`
}

type Response struct {
	response.Response
	Jwt string `json:"jwt"`
}

type LoginService interface {
	Login(ctx context.Context, email string, password string, appID int32) (string, error)
}

func New(log *slog.Logger, loginService LoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login"

		log := log.With(
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

		jwt, err := loginService.Login(context.Background(), req.Email, req.Password, req.AppID)
		if err != nil {
			if errors.Is(err, grpc.ErrUserNotFound) {
				log.Info("the user does not exist", slog.String("email", req.Email))

				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, response.Error("the user does not exist"))

				return
			}

			log.Info("failed to login user", slog.String("email", req.Email))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to login user"))

			return
		}

		log.Info("user loggined")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{
			Jwt: jwt,
		})
	}
}
