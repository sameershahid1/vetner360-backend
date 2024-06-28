package routes

import (
	"vetner360-backend/controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleMobileRoutes(router chi.Router) {
	// router.Use(custom_middleware.ValidateJsonFormat)
	router.Use(middleware.CleanPath)

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)
		protectedRoute.Get("/doctors", controller.GetPetOwners)

		protectedRoute.Route("/profile", func(moduleRoute chi.Router) {
			moduleRoute.Get("/{id}", controller.GetProfile)
			moduleRoute.Patch("/{id}", controller.GetProfile)
		})

		protectedRoute.Route("/pet", func(moduleRoute chi.Router) {
			// moduleRoute.Get("/detail/{userId}/{id}", controller.GetPetDetail)
			moduleRoute.Post("/my-pet/{userId}", controller.GetMyPetList)
			moduleRoute.Post("/", controller.PostPet)
			moduleRoute.Patch("/{id}", controller.PatchPet)
			moduleRoute.Delete("/{userId}/{id}", controller.DeletePet)

			moduleRoute.Route("/activity", func(subModuleRoute chi.Router) {
				subModuleRoute.Post("/list/{petId}", controller.GetActivityList)
				subModuleRoute.Post("/", controller.PostActivity)
				subModuleRoute.Patch("/{id}", controller.PatchActivity)
				subModuleRoute.Delete("/{petId}/{id}", controller.DeleteActivity)
			})
		})

		protectedRoute.Route("/doctor", func(moduleRoute chi.Router) {
			moduleRoute.Get("/nearest", controller.GetNearestDoctors)
		})

		protectedRoute.Route("/chat", func(moduleRoute chi.Router) {
			moduleRoute.Post("/participant/{roomId}", controller.GetChatParticipant)
			moduleRoute.Post("/chat-participant/add", controller.AddParticipant)
			moduleRoute.Get("/messages/{roomId}", controller.GetChatMessages)
		})

		// protectedRoute.Route("/health-report", func(moduleRoute chi.Router) {
		// 	moduleRoute.Post("/my-pet/{id}", controller.GetPetActivity)
		// 	moduleRoute.Post("/", controller.PostPet)
		// 	moduleRoute.Patch("/{id}", controller.PatchPet)
		// 	moduleRoute.Delete("/{id}", controller.DeletePet)
		// })

	})

	router.Group(func(publicRoute chi.Router) {
		publicRoute.Post("/login", controller.SignIn)
		publicRoute.Post("/user-registration", controller.PetOwnerORGuestRegistration)
		publicRoute.Post("/doctor-registration", controller.DoctorRegistration)
	})
}
