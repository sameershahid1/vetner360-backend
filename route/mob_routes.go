package routes

import (
	"vetner360-backend/controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleMobileRoutes(router chi.Router) {

	router.Use(custom_middleware.ValidateJsonFormat)
	router.Use(middleware.CleanPath)

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)
		protectedRoute.Get("/doctors", controller.GetPetOwners)
		protectedRoute.Route("/profile", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetProfile)
			moduleRoute.Patch("/{id}", controller.GetProfile)
		})

		protectedRoute.Route("/pet", func(moduleRoute chi.Router) {
			moduleRoute.Post("/my-pet/{id}", controller.GetMyPetList)
			moduleRoute.Get("/detail/{userId}/{id}", controller.GetPetDetail)
			moduleRoute.Post("/", controller.PostPet)
			moduleRoute.Patch("/{id}", controller.PatchPet)
			moduleRoute.Delete("/{id}", controller.DeletePet)
		})
	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.SignIn)
		publicRoute.Post("/user-registration", controller.PetOwnerORGuestRegistration)
		publicRoute.Post("/doctor-registration", controller.DoctorRegistration)
	})
}
