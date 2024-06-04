package controller

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

var SocketServer = socketio.NewServer(&engineio.Options{
	Transports: []transport.Transport{
		&polling.Transport{
			CheckOrigin: allowOriginFunc,
		},
		&websocket.Transport{
			CheckOrigin: allowOriginFunc,
		},
	},
})

func SocketConnection(s socketio.Conn) error {
	s.SetContext("")
	log.Println("connected:", s.ID())
	return nil
}

func EventMessage(s socketio.Conn, msg string) string {

	SocketServer.BroadcastToRoom("/", "ab1", "reply", msg)
	return "recv " + msg
}

func JoinRoom(s socketio.Conn, msg string) {
	//Check if Participants exists, if exists than it will create room, other wise it will ignore it
	// if 1 {
	// }
	s.Join("ab1")
	s.Emit("reply", "Have Joined "+msg)
}

func LeaveRoom(s socketio.Conn) string {
	s.Leave("ab1")
	s.Emit("reply", "Leaving the room")
	return "Leaving the room"
}

func SocketError(s socketio.Conn, e error) {
	log.Println("meet error:", e)
}

func SocketDisconnect(s socketio.Conn, reason string) {
	s.LeaveAll()
	log.Println("closed", reason)
}

func GetChatMessage(response http.ResponseWriter, request *http.Request) {

	response.Write([]byte("Messages"))
}
