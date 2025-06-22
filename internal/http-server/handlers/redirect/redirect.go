package redirect

import (
	"log/slog"
	"net/http"

	"github.com/arvinloc/url-shortener/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("is empty"))

			return
		}

		foundURL, err := urlGetter.GetURL(alias)

		if err != nil {
			log.Error("url not found", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("url not found"))
		}

		log.Info("alias found")

		http.Redirect(w, r, foundURL, http.StatusFound)
	}
}
