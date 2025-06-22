package delete

import (
	"log/slog"
	"net/http"

	"github.com/arvinloc/url-shortener/internal/http-server/handlers/payload"
	"github.com/arvinloc/url-shortener/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest) // Добавьте эту строку
			render.JSON(w, r, response.Error("is empty"))
			return
		}
		err := urlDeleter.DeleteURL(alias)

		if err != nil {
			log.Error("failed to delete url", slog.String("err", err.Error()))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("url not found"))
			return
		}

		log.Info("url deleted")

		render.JSON(w, r, payload.Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
