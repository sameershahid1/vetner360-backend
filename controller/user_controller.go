package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	// var result struct {
	// 	ID       primitive.ObjectID `bson:"_id"`
	// 	Name     string
	// 	Password string
	// 	Email    string
	// }
	// collection := database.MongoDB.Database("vetner360").Collection("users")
	// filter := bson.D{{"name", "Sameer"}}
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// err := collection.FindOne(ctx, filter).Decode(&result)
	// if err == mongo.ErrNoDocuments {
	// 	fmt.Println("record does not exist")
	// } else if err != nil {
	// 	log.Fatal(err)
	// }

	// var response postResposne = postResposne{Status: true, Message: "Successfully "}
	// jsonData, err := json.Marshal(response)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// w.Write(jsonData)

	jsonData, err := json.Marshal(Response{Status: true, Message: "Hi, admin"})
	if err != nil {
		log.Fatal(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Write(jsonData)
}
func Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte("The king is back"))
}
