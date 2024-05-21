package mongodb

import (
	"context"
	"time"
	"vetner360-backend/database"
	data_type "vetner360-backend/utils/type"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database string
var Collection string

func GetAll[T data_type.RecordType](filter *bson.M, opts *options.FindOptions) ([]T, error) {
	var records []T
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, *filter, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	errorCur := cur.All(context.Background(), &records)
	if errorCur != nil {
		return nil, errorCur
	}

	return records, nil
}

func GetOne[T data_type.RecordType](filter bson.M) (*T, error) {
	var record T
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()
	// bson.M{"_id": id}
	errCur := collection.FindOne(ctx, filter).Decode(&record)
	if errCur != nil {
		return nil, errCur
	}

	return &record, nil
}

func Post[T data_type.RecordType](data bson.M) (interface{}, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	response, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return response.InsertedID, nil
}

func Patch[T data_type.RecordType](filter bson.M, updatedData bson.M) (*mongo.UpdateResult, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	record, errCur := collection.UpdateOne(ctx, filter, bson.M{"$set": updatedData})
	if errCur != nil {
		return nil, errCur
	}

	return record, nil
}

func Delete[T data_type.RecordType](filter bson.M) (*mongo.DeleteResult, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	result, errCur := collection.DeleteOne(ctx, filter)
	if errCur != nil {
		return nil, errCur
	}

	return result, nil
}
