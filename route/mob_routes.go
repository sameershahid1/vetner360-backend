package routes

import (
	"vetner360-backend/controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
)

func HandleMobileRoutes(router chi.Router) {
	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)
		protectedRoute.Get("/doctors", controller.GetPetOwners)
		// protectedRoute.HandleFunc("/chat", controller.SocketHandling)

		// protectedRoute.Mount("/socket.io", server)
	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.SignIn)
		publicRoute.Post("/user-registration", controller.PetOwnerORGuestRegistration)
		publicRoute.Post("/doctor-registration", controller.DoctorRegistration)
	})
}
