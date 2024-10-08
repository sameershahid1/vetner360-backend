package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"vetner360-backend/controller"
	"vetner360-backend/controller/mobile_controller"
	"vetner360-backend/database"
	routes "vetner360-backend/route"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	url := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.ConnectWithMongoDB(url)
	defer database.DisconnectWithMongodb()

	// for collectionName, columnList := range static_data.IndexCollection {
	// 	for column := range columnList {
	// 		collectionAttribute := static_data.IndexCollectionAttribute[collectionName]
	// 		database.IndexingCollection(collectionName, columnList[column], collectionAttribute[columnList[column]])
	// 	}
	// }

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
	router.Use(middleware.Timeout(time.Second * 60))

	mobile_controller.SocketServer.OnConnect("/", mobile_controller.SocketConnection)
	mobile_controller.SocketServer.OnEvent("/", "message", mobile_controller.EventMessage)
	mobile_controller.SocketServer.OnEvent("/", "join-room", mobile_controller.JoinRoom)
	mobile_controller.SocketServer.OnEvent("/", "leave-room", mobile_controller.LeaveRoom)
	mobile_controller.SocketServer.OnError("/", mobile_controller.SocketError)
	mobile_controller.SocketServer.OnDisconnect("/", mobile_controller.SocketDisconnect)

	go func() {
		if err := mobile_controller.SocketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer mobile_controller.SocketServer.Close()

	router.Handle("/socket.io/", mobile_controller.SocketServer)
	router.Mount("/debug", middleware.Profiler())
	router.Route("/web/api", routes.HandleWebRoutes)
	router.Route("/mobile/api", routes.HandleMobileRoutes)
	router.NotFound(controller.RouteDoesExists)
	router.MethodNotAllowed(controller.MethodNotExists)
	router.Handle("/*", http.FileServer(http.Dir("./public")))

	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
