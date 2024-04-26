package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/pkg/api/apiv1"
)

func Router(handler apiv1.ServerInterface) http.Handler {
	router := chi.NewRouter()
	router.Mount(
		"/api", apiv1.HandlerWithOptions(
			handler, apiv1.ChiServerOptions{
				BaseURL: "/v1",
				ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					slog.Error("handle error", slog.String("err", err.Error()))
				},
			},
		),
	)
	return router
}
