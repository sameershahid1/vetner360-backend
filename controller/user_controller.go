package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"vetner360-backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type postResposne struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID       primitive.ObjectID `bson:"_id"`
		Name     string
		Password string
		Email    string
	}
	collection := database.MongoDB.Database("vetner360").Collection("users")
	filter := bson.D{{"name", "Sameer"}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
	} else if err != nil {
		log.Fatal(err)
	}

	var response postResposne = postResposne{Status: true, Message: "Successfully "}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(jsonData)
}
