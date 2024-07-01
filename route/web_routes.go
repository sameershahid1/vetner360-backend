package routes

import (
	"vetner360-backend/controller"
	"vetner360-backend/controller/mobile_controller"
	"vetner360-backend/controller/web_controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleWebRoutes(router chi.Router) {
	router.Use(middleware.CleanPath)
	// router.Use(custom_middleware.ValidateJsonFormat)

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)

		protectedRoute.Route("/user", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list", web_controller.GetUser)
			moduleRoute.Post("/", web_controller.PostUser)
			moduleRoute.Patch("/{id}", web_controller.PatchUser)
			moduleRoute.Delete("/{id}", web_controller.DeleteUser)
		})

		protectedRoute.Route("/doctor", func(moduleRoute chi.Router) {
			moduleRoute.Get("/", controller.GetDoctor)
			moduleRoute.Post("/", controller.PostDoctor)
			moduleRoute.Patch("/{id}", controller.PatchDoctor)
			moduleRoute.Patch("/status/{id}", controller.PatchDoctor)
			moduleRoute.Delete("/{id}", controller.DeleteDoctor)
		})

		protectedRoute.Route("/pet", func(moduleRoute chi.Router) {
			moduleRoute.Post("/my-pet/{userId}", controller.GetMyPetList)
			moduleRoute.Post("/", controller.PostPet)
			moduleRoute.Patch("/{id}", controller.PatchPet)
			moduleRoute.Delete("/{userId}/{id}", controller.DeletePet)

			moduleRoute.Route("/activity", func(subModuleRoute chi.Router) {
				subModuleRoute.Post("/list/{petId}", mobile_controller.GetActivityList)
				subModuleRoute.Post("/", mobile_controller.PostActivity)
				subModuleRoute.Patch("/{id}", mobile_controller.PatchActivity)
				subModuleRoute.Delete("/{petId}/{id}", mobile_controller.DeleteActivity)
			})
		})

		protectedRoute.Route("/role", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list", web_controller.GetRoles)
			moduleRoute.Post("/", web_controller.PostRoleOwner)
			moduleRoute.Patch("/{id}", web_controller.GetRoles)
			moduleRoute.Delete("/{id}", web_controller.GetRoles)
		})

		protectedRoute.Route("/contact-message", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list", web_controller.GetContactMessage)
			moduleRoute.Post("/", web_controller.PostContactMessage)
			moduleRoute.Patch("/{id}", web_controller.PatchContactMessage)
			moduleRoute.Delete("/{id}", web_controller.DeleteContactMessage)
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
