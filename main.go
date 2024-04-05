package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"vetner360-backend/controller"
	"vetner360-backend/database"
	routes "vetner360-backend/route"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	url := os.Getenv("MONGODB_URI")
	database.ConnectWithMongoDB(url)
	defer database.DisconnectWithMongodb()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	router.Route("/web/api", routes.HandleWebRoutes)
	router.Route("/mobile/api", routes.HandleMobileRoutes)

	router.NotFound(controller.RouteDoesExists)
	router.MethodNotAllowed(controller.MethodNotExists)

	http.ListenAndServe(":8080", router)
}
