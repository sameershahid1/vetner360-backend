package routes

import (
	"vetner360-backend/controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleWebRoutes(router chi.Router) {
	router.Use(middleware.CleanPath)
	router.Use(custom_middleware.ValidateJsonFormat)

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
			moduleRoute.Get("/", controller.GetProfile)
			moduleRoute.Post("/", controller.GetProfile)
			moduleRoute.Patch("/{id}", controller.GetProfile)
			moduleRoute.Delete("/{id}", controller.GetProfile)
		})

		protectedRoute.Route("/Guests", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetGuests)
			moduleRoute.Post("/", controller.GetGuests)
			moduleRoute.Patch("/{id}", controller.GetGuests)
			moduleRoute.Delete("/{id}", controller.GetGuests)
		})

		protectedRoute.Route("/role", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list", controller.GetRoles)
			moduleRoute.Post("/", controller.PostRoleOwner)
			moduleRoute.Patch("/{id}", controller.GetRoles)
			moduleRoute.Delete("/{id}", controller.GetRoles)
		})

		// protectedRoute.Route("/permission", func(moduleRoute chi.Router) {
		// 	moduleRoute.Get("/", controller.GetPermissions)
		// 	moduleRoute.Post("/", controller.GetPermissions)
		// 	moduleRoute.Patch("/{id}", controller.GetPermissions)
		// 	moduleRoute.Delete("/{id}", controller.GetPermissions)
		// })

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
