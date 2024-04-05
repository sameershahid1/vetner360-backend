package routes

import (
	"vetner360-backend/controller"
	"vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
)

func HandleMobileRoutes(router chi.Router) {
	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(middleware.MobileVerifyJWT)
		protectedRoute.Get("/doctors", controller.GetPetOwners)
	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.MobileSignIn)
	})
}
