package routes

import (
	"vetner360-backend/controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
)

func HandleWebRoutes(router chi.Router) {

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)

		protectedRoute.Route("/pet-owner", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list", controller.GetPetOwners)
			moduleRoute.Post("/", controller.PostPetOwner)
			moduleRoute.Patch("/{id}", controller.PatchPetOwner)
			moduleRoute.Delete("/{id}", controller.DeletePetOwner)
		})

		protectedRoute.Route("/doctor", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetDoctors)
			moduleRoute.Post("/", controller.GetDoctors)
			moduleRoute.Patch("/{id}", controller.GetDoctors)
			moduleRoute.Delete("/{id}", controller.GetDoctors)
		})

		protectedRoute.Route("/pets", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetPets)
			moduleRoute.Post("/", controller.GetPets)
			moduleRoute.Patch("/{id}", controller.GetPets)
			moduleRoute.Delete("/{id}", controller.GetPets)
		})

		protectedRoute.Route("/Guests", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetGuests)
			moduleRoute.Post("/", controller.GetGuests)
			moduleRoute.Patch("/{id}", controller.GetGuests)
			moduleRoute.Delete("/{id}", controller.GetGuests)
		})

		protectedRoute.Route("/role", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetRoles)
			moduleRoute.Post("/", controller.GetRoles)
			moduleRoute.Patch("/{id}", controller.GetRoles)
			moduleRoute.Delete("/{id}", controller.GetRoles)
		})

		protectedRoute.Route("/permission", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetPermissions)
			moduleRoute.Post("/", controller.GetPermissions)
			moduleRoute.Patch("/{id}", controller.GetPermissions)
			moduleRoute.Delete("/{id}", controller.GetPermissions)
		})

		protectedRoute.Route("/contact-message", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetContactMessages)
		})

		protectedRoute.Route("/profile", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetProfile)
			moduleRoute.Patch("/{id}", controller.GetProfile)
		})
	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.SignIn)
	})
}
