package save

import (
	"log/slog"
	"net/http"

	"github.com/arvinloc/url-shortener/internal/http-server/handlers/payload"
	"github.com/arvinloc/url-shortener/internal/lib/api/response"
	"github.com/arvinloc/url-shortener/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(6)
		}

		err = urlSaver.SaveURL(req.URL, alias)

		if err != nil {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		log.Info("url added")

		render.JSON(w, r, payload.Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
