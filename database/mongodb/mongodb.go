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

var Database string = "vetner360"

func GetAll[T data_type.RecordType](filter *bson.M, opts *options.FindOptions, Collection string) ([]T, error) {
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

func GetAllUsingPipeline[T data_type.RecordType](pipeline mongo.Pipeline, opts *options.AggregateOptions, Collection string) ([]T, error) {
	var records []T
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetOne[T data_type.RecordType](filter bson.M, opts *options.FindOneOptions, Collection string) (*T, error) {
	var record T
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	errCur := collection.FindOne(ctx, filter, opts).Decode(&record)
	if errCur != nil {
		return nil, errCur
	}

	return &record, nil
}

func Post[T data_type.RecordType](data bson.M, Collection string) (interface{}, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	response, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return response.InsertedID, nil
}

func Patch[T data_type.RecordType](filter bson.M, updatedData bson.M, Collection string) (*mongo.UpdateResult, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	record, errCur := collection.UpdateOne(ctx, filter, bson.M{"$set": updatedData})
	if errCur != nil {
		return nil, errCur
	}

	return record, nil
}

func Delete[T data_type.RecordType](filter bson.M, Collection string) (*mongo.DeleteResult, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	result, errCur := collection.DeleteOne(ctx, filter)
	if errCur != nil {
		return nil, errCur
	}

	return result, nil
}

func TotalDocs[T data_type.RecordType](filter *bson.M, Collection string) (int64, error) {
	collection := database.MongoDB.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout*time.Second)
	defer cancel()

	result, errCur := collection.CountDocuments(ctx, filter)
	if errCur != nil {
		return 0, errCur
	}

	return result, nil
}
