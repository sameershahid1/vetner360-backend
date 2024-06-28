package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDoctors(response http.ResponseWriter, request *http.Request) {
	requestResponse := data_type.Response[model.Doctor]{Status: true, Message: "Hi, admin"}
	jsonData, _ := json.Marshal(requestResponse)
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func GetNearestDoctors(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	if query["latitude"] == nil {
		helping.InternalServerError(response, errors.New("missing latitude query"))
		return
	}
	if query["longitude"] == nil {
		helping.InternalServerError(response, errors.New("missing longitude query"))
		return
	}

	latitude, err := strconv.ParseFloat(query["latitude"][0], 32)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}
	longitude, err := strconv.ParseFloat(query["longitude"][0], 32)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	mongoDBHQ := bson.M{
		"type":        "Point",
		"coordinates": []float64{longitude, latitude}, // Use float64 for higher precision
	}
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry":    mongoDBHQ,
				"$minDistance": 0,
				"$maxDistance": 10000,
			},
		},
	}
	opts := options.FindOptions{}

	records, err := mongodb.GetAll[model.Doctor](&filter, &opts, "users")

	if err != nil {
		print(err.Error())
		helping.InternalServerError(response, err)
		return
	}

	if records == nil {
		records = []model.Doctor{}
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully loaded nearest doctors", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
