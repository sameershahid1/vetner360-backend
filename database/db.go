package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client
var ctxDB *context.Context
var errDb *error

const QueryTimeout = 10

func ConnectWithMongoDB(envFileName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	initialInterval := 1 * time.Second
	maxInterval := 10 * time.Second
	for {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(envFileName))
		if err == nil {
			MongoDB = client
			ctxDB = &ctx
			errDb = &err
			fmt.Println("Successfully connected with database")
			break
		}

		select {
		case <-ctx.Done():
			fmt.Println("database connection timeout")
		case <-time.After(initialInterval):
			fmt.Println("Error connecting to MongoDB:", err)
			fmt.Println("Retrying in", initialInterval)
			if initialInterval*2 <= maxInterval {
				initialInterval = initialInterval * 2
			} else {
				initialInterval = maxInterval
			}
		}

	}

}

func IndexingCollection(collectionName string, attribute string, value interface{}) {
	if MongoDB == nil {
		return
	}
	collection := MongoDB.Database("vetner360").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: attribute, Value: value}},
		Options: options.Index().SetBackground(true).SetName(attribute + "_index"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(collectionName + " attribute " + attribute + " is indexed")
	}
}

func DisconnectWithMongodb() {
	if MongoDB != nil {
		return
	}
	*errDb = MongoDB.Disconnect(*ctxDB)
	if *errDb != nil {
		fmt.Println("Disconnected with database")
		recover()
	}
}
