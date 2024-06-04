package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"vetner360-backend/controller"
	"vetner360-backend/database"
	routes "vetner360-backend/route"
	static_data "vetner360-backend/utils/data"

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
	for i := 0; i < len(static_data.IndexCollection); i++ {
		collectionName := static_data.IndexCollection[i]
		collectionAttribute := static_data.IndexCollectionAttribute[collectionName]
		database.IndexingCollection(collectionName, collectionAttribute)
	}
	defer database.DisconnectWithMongodb()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(middleware.Timeout(time.Second * 60))

	controller.SocketServer.OnConnect("/", controller.SocketConnection)
	controller.SocketServer.OnEvent("/", "msg", controller.EventMessage)
	controller.SocketServer.OnEvent("/", "join-room", controller.JoinRoom)
	controller.SocketServer.OnEvent("/", "leave-room", controller.LeaveRoom)
	controller.SocketServer.OnError("/", controller.SocketError)
	controller.SocketServer.OnDisconnect("/", controller.SocketDisconnect)

	go func() {
		if err := controller.SocketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer controller.SocketServer.Close()

	router.Handle("/socket.io/", controller.SocketServer)
	router.Mount("/debug", middleware.Profiler())
	router.Route("/web/api", routes.HandleWebRoutes)
	router.Route("/mobile/api", routes.HandleMobileRoutes)
	router.NotFound(controller.RouteDoesExists)
	router.MethodNotAllowed(controller.MethodNotExists)

	http.ListenAndServe(":8080", router)
}
