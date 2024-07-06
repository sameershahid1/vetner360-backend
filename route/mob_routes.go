package routes

import (
	"vetner360-backend/controller"
	"vetner360-backend/controller/mobile_controller"
	custom_middleware "vetner360-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleMobileRoutes(router chi.Router) {
	// router.Use(custom_middleware.ValidateJsonFormat)
	router.Use(middleware.CleanPath)

	router.Group(func(protectedRoute chi.Router) {
		protectedRoute.Use(custom_middleware.VerifyJWTMiddleware)

		protectedRoute.Route("/profile", func(moduleRoute chi.Router) {
			moduleRoute.Get("/{id}", controller.GetProfile)
			moduleRoute.Patch("/user/{id}", controller.UpdateUserProfile)
			moduleRoute.Patch("/doctor/{id}", controller.UpdateDoctorProfile)
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

		protectedRoute.Route("/doctor", func(moduleRoute chi.Router) {
			moduleRoute.Get("/nearest", controller.GetNearestDoctors)
			moduleRoute.Get("/search-location", controller.GetNearestDoctors)
		})

		protectedRoute.Route("/chat", func(moduleRoute chi.Router) {
			moduleRoute.Post("/participant/{userId}", mobile_controller.GetChatParticipant)
			moduleRoute.Post("/chat-participant/add", mobile_controller.AddParticipant)
			moduleRoute.Get("/messages/{roomId}", mobile_controller.GetChatMessages)
			moduleRoute.Get("/messages/latest/{roomId}", mobile_controller.GetLatestMessage)
		})

		protectedRoute.Route("/appointment", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list/", mobile_controller.GetActivityList)
			moduleRoute.Post("/", mobile_controller.PostActivity)
			moduleRoute.Patch("/{id}", mobile_controller.PatchActivity)
			moduleRoute.Delete("/{petId}/{id}", mobile_controller.DeleteActivity)
		})

		protectedRoute.Route("/sell", func(moduleRoute chi.Router) {
			moduleRoute.Post("/list/{userId}", mobile_controller.GetMySellPets)
			moduleRoute.Post("/", mobile_controller.PostMyPetToSell)
			moduleRoute.Patch("/{id}", mobile_controller.PatchMyPetOnSell)
			moduleRoute.Delete("/{id}", mobile_controller.DeleteMyPetOnSell)
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
		publicRoute.Post("/user-registration", controller.UserRegistration)
		publicRoute.Post("/doctor-registration", controller.DoctorRegistration)

		publicRoute.Post("/doctor/latest/clinic", mobile_controller.GetLatestDoctorClinic)
		publicRoute.Post("/pet/latest/{type}", mobile_controller.GetLatestPetByType)

	})
}
