package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	methodController "github.com/vetner360-backend/controller"
	userController "github.com/vetner360-backend/controller"
)

func main() {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)

	//Requests Path
	router.Get("/", userController.GetUser)

	router.NotFound(methodController.RouteDoesExists)
	router.MethodNotAllowed(userController.MethodNotExists)

	http.ListenAndServe(":8080", router)
}
