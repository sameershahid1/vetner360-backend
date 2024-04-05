package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPetOwners(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"
	records, err := mongodb.GetAll[model.User]()
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request", Records: records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostPetOwners(response http.ResponseWriter, request *http.Request) {
	var newRecord = bson.M{"firstName": "sameer", "lastName": "shahid", "email": "sam123@gmail.com", "phoneNo": "0321456789", "password": "123456789", "created_at": time.Now()}
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"
	_, err := mongodb.Post(newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
