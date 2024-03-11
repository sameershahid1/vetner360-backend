package routes

import (
	"vetner360-backend/controller"

	"github.com/go-chi/chi/v5"
)

func HandleWebRoutes(router chi.Router) {
	router.Get("/", controller.GetUser)
}
