package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client
var ctxDB *context.Context
var errDb *error

const QueryTimeout = 10

func ConnectWithMongoDB(envFileName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(envFileName))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	MongoDB = client
	ctxDB = &ctx
	errDb = &err
	fmt.Println("Successfully connected with database")
}

func DisconnectWithMongodb() {
	if *errDb = MongoDB.Disconnect(*ctxDB); *errDb != nil {
		fmt.Println("Disconnected with database")
		panic(errDb)
	}
}
