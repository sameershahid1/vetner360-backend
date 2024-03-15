package routes

import (
	"vetner360-backend/controller"
	"vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
)

func HandleWebRoutes(router chi.Router) {

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(middleware.VerifyJWT)
		protectedRoute.Get("/admin", controller.GetUser)
	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.Signin)
	})
}
