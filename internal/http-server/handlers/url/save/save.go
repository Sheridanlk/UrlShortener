package save

import (
	"UrlShortener/internal/lib/api/response"
	"UrlShortener/internal/lib/authctx"
	"UrlShortener/internal/lib/random"
	"UrlShortener/internal/logger/sl"
	"UrlShortener/internal/storage"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

//go:generate mockery
type URLSaver interface {
	SaveURL(urlToSave, alias string, userID int64) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		UserIdnCtx := r.Context().Value(authctx.UserIDKey).(int64)

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

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomAlias(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias, UserIdnCtx)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("alias already exists", slog.String("url", req.URL))

			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Error("alias already exists"))

			return
		}
		if err != nil {
			log.Info("failed to add url ", slog.String("url", req.URL))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("ulr added", slog.Int64("id", id))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{
			Alias: alias,
		})
	}
}
