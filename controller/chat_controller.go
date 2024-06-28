package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	fmt.Println("connected:", s.ID())
	return nil
}

func EventMessage(s socketio.Conn, msg data_type.MessageBody) {
	id := uuid.New()
	sendMessage := map[string]interface{}{"senderId": msg.SenderId,
		"roomId":     msg.RoomId,
		"content":    msg.Content,
		"type":       msg.Type,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	var newRecord = bson.M{
		"senderId":   msg.SenderId,
		"roomId":     msg.RoomId,
		"content":    msg.Content,
		"type":       msg.Type,
		"token":      id.String(),
		"created_at": time.Now(),
	}

	_, err := mongodb.Post[model.Message](newRecord, "messages")
	if err != nil {
		fmt.Println(err.Error())
	}

	SocketServer.BroadcastToRoom("/", msg.RoomId, "received-message", sendMessage)
}

func JoinRoom(s socketio.Conn, room string) {
	roomList := s.Rooms()
	isRoom := true
	response := map[string]interface{}{"status": true, "response": "Joined the Chat"}
	for x := range roomList {
		if roomList[x] == room {
			isRoom = false
			response["status"] = false
			response["response"] = "Already join the chat"
		}
	}
	if isRoom {
		s.Join(room)
	}
	SocketServer.BroadcastToRoom("/", room, "joined-room", response)
}

func LeaveRoom(s socketio.Conn, room string) {
	// s.Leave(room)
}

func SocketError(s socketio.Conn, e error) {
	log.Println("meet error:", e)
}

func SocketDisconnect(s socketio.Conn, reason string) {
	s.LeaveAll()
	fmt.Println("closed", reason)
}

func GetChatParticipant(response http.ResponseWriter, request *http.Request) {

	var requestBody data_type.PaginationType[model.Participant]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var roomId = chi.URLParam(request, "roomId")
	var filter = bson.M{"roomId": roomId}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Participant](&filter, &opts, "participants")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	if records == nil {
		records = []model.Participant{}
	}

	var requestResponse = data_type.Response[model.Participant]{Status: true, Message: "Successfully loaded participant Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetChatMessages(response http.ResponseWriter, request *http.Request) {
	var roomId = chi.URLParam(request, "roomId")
	var filter = bson.M{"roomId": roomId}
	opts := options.FindOptions{}
	opts.Sort = bson.D{{"created_at", -1}}

	records, err := mongodb.GetAll[model.Message](&filter, &opts, "messages")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	if records == nil {
		records = []model.Message{}
	}

	var requestResponse = data_type.Response[model.Message]{Status: true, Message: "Successfully loaded messages", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func AddParticipant(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.ParticipantType
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(requestBody)
	if err != nil {
		errorMessageList := strings.Split(err.Error(), "\n")
		errorMessage := strings.Split(errorMessageList[0], "Error:")
		response.WriteHeader(http.StatusBadRequest)
		jsonErrorMessage, err := helping.JsonEncode(errorMessage[1])
		if err != nil {
			response.Write([]byte("Internal side error"))
		}
		response.Write(jsonErrorMessage)
		return
	}

	isSame, _ := mongodb.GetOne[model.Participant](bson.M{"userId": requestBody.UserId, "roomId": requestBody.RoomId}, "participants")
	if isSame != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Participants already exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var newRecord = bson.M{
		"userId":     requestBody.UserId,
		"roomId":     requestBody.RoomId,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.Participant](newRecord, "participants")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Participant]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
