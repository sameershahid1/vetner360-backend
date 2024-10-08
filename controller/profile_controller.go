package controller

import (
	"encoding/json"
	"net/http"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.User]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var userId = chi.URLParam(request, "id")
	var filter = bson.M{"token": userId}
	opts := options.FindOneOptions{}
	records, err := mongodb.GetOne[model.User](filter, &opts, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request", Data: records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func UpdateUserProfile(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")

	var requestBody data_type.PetOwnerRequestType
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.User](filter, &opts, "users")
	if isSameUser == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("User does not exists")
		if err != nil {
			helping.InternalServerError(response, err, http.StatusInternalServerError)
			return
		}
		response.Write(jsonResponse)
		return
	}

	cost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), cost)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var updateRecord = bson.M{
		"firstName": requestBody.FirstName,
		"lastName":  requestBody.LastName,
		"email":     requestBody.Email,
		"phoneNo":   requestBody.PhoneNo,
		"password":  string(bytes),
	}

	_, err = mongodb.Patch[model.User](filter, updateRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully updated pet"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func UpdateDoctorProfile(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var requestBody data_type.DoctorRequestType
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.Doctor](filter, &opts, "users")
	if isSameUser == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("User does not exists")
		if err != nil {
			helping.InternalServerError(response, err, http.StatusInternalServerError)
			return
		}
		response.Write(jsonResponse)
		return
	}

	cost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), cost)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var location = model.Location{
		Type:        "Point",
		Coordinates: []float64{requestBody.Longitude, requestBody.Latitude},
	}

	var updateRecord = bson.M{
		"firstName":    requestBody.FirstName,
		"lastName":     requestBody.LastName,
		"email":        requestBody.Email,
		"phoneNo":      requestBody.PhoneNo,
		"password":     string(bytes),
		"fatherName":   requestBody.FatherName,
		"registration": requestBody.Registration,
		"clinicName":   requestBody.ClinicName,
		"location":     location,
	}

	_, err = mongodb.Patch[model.Doctor](filter, updateRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully updated pet"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
