package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"

	var requestBody data_type.PaginationType[model.User]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var userId = chi.URLParam(request, "id")
	var filter = bson.M{"token": userId}

	records, err := mongodb.GetOne[model.User](filter)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request", Data: records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func UpdateProfile(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"

	var requestBody data_type.PetOwnerRequestType
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

	var filter = bson.M{"token": id}
	isSameUser, _ := mongodb.GetOne[model.User](filter)
	if isSameUser == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("User does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	cost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), cost)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var updateRecord = bson.M{
		"firstName": requestBody.FirstName,
		"lastName":  requestBody.LastName,
		"email":     requestBody.Email,
		"phoneNo":   requestBody.PhoneNo,
		"password":  string(bytes),
	}

	_, err = mongodb.Patch[model.User](filter, updateRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully updated pet"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
