package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandleMobileRoutes(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fuck mobile"))
	})
}
