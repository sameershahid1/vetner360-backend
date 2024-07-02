package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func GetDoctor(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.Doctor]
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

	var filter = bson.M{"roleId": "665cec7fc6206b06eddaacca"}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Doctor](&filter, &opts, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	total, err := mongodb.TotalDocs[model.User](&filter, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully Completed Request", Records: &records, Count: &total}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostDoctor(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
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
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.Doctor](bson.M{"email": requestBody.Email}, &opts, "users")
	if isSameUser != nil {
		helping.InternalServerError(response, errors.New("error: doctor already exists"), http.StatusBadRequest)
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

	var newRecord = bson.M{
		"firstName":    requestBody.FirstName,
		"lastName":     requestBody.LastName,
		"email":        requestBody.Email,
		"phoneNo":      requestBody.PhoneNo,
		"password":     string(bytes),
		"fatherName":   requestBody.FatherName,
		"registration": requestBody.Registration,
		"clinicName":   requestBody.ClinicName,
		"experience":   requestBody.Experience,
		"bio":          requestBody.Bio,
		"location":     location,
		"token":        id.String(),
		"roleId":       "665cec7fc6206b06eddaacca",
		"created_at":   time.Now(),
	}
	_, err = mongodb.Post[model.Doctor](newRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully Register Doctor"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchDoctor(response http.ResponseWriter, request *http.Request) {
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
		helping.InternalServerError(response, errors.New("error: doctor does not exists"), http.StatusBadRequest)
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

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully updated doctor"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchDoctorStatus(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var requestBody data_type.DoctorStatusRequestType
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
		helping.InternalServerError(response, errors.New("error: doctor does not exists"), http.StatusBadRequest)
		return
	}

	var updateRecord = bson.M{
		"status": requestBody.Status,
	}

	_, err = mongodb.Patch[model.Doctor](filter, updateRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully updated doctor status"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteDoctor(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.Doctor](filter, &opts, "users")
	if isSameUser == nil {
		helping.InternalServerError(response, errors.New("error: user does not exists"), http.StatusBadRequest)
		return
	}

	_, err := mongodb.Delete[model.Doctor](filter, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully deleted pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetNearestDoctors(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	opts := options.FindOptions{}

	if query["latitude"] == nil {
		helping.InternalServerError(response, errors.New("missing latitude query"), http.StatusBadRequest)
		return
	}
	if query["longitude"] == nil {
		helping.InternalServerError(response, errors.New("missing longitude query"), http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(query["latitude"][0], 32)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}
	longitude, err := strconv.ParseFloat(query["longitude"][0], 32)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	mongoDBHQ := bson.M{
		"type":        "Point",
		"coordinates": []float64{longitude, latitude},
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

	records, err := mongodb.GetAll[model.Doctor](&filter, &opts, "users")

	if err != nil {
		print(err.Error())
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	if records == nil {
		records = []model.Doctor{}
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully loaded nearest doctors", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetLocations(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Locations"))
}
