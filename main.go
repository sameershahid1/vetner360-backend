package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"vetner360-backend/database"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

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
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	// router.Use(middleware.CleanPath)
	// router.Use(httprate.LimitByIP(100, 1*time.Minute))
	// router.Use(middleware.Timeout(time.Second * 60))
	// router.Use(custom_middleware.ValidateJsonFormat)

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		log.Println("chat:", msg)
		fmt.Println(s.Rooms())
		server.BroadcastToRoom("/", "ab1", "reply", "FUCK YOU")
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, msg string) {
		s.Join("ab1")
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "leave-room", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Handle("/socket.io/", server)

	// router.Mount("/debug", middleware.Profiler())

	// router.Route("/web/api", routes.HandleWebRoutes)
	// router.Route("/mobile/api", routes.HandleMobileRoutes)

	// router.NotFound(controller.RouteDoesExists)
	// router.MethodNotAllowed(controller.MethodNotExists)

	http.ListenAndServe(":8080", router)
}
