package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func GetPetOwners(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"

	var requestBody data_type.PaginationType[model.User]
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
		helping.JsonEncode(errorMessage[1])
		return
	}

	var filter = bson.M{}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.User](&filter, &opts)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostPetOwner(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"
	id := uuid.New()
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
		helping.JsonEncode(errorMessage[1])
		return
	}

	isSameUser, _ := mongodb.GetOne[model.User](bson.M{"email": requestBody.Email})
	if isSameUser != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("User already exists")
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

	var newRecord = bson.M{
		"firstName":  requestBody.FirstName,
		"lastName":   requestBody.LastName,
		"email":      requestBody.Email,
		"phoneNo":    requestBody.PhoneNo,
		"password":   string(bytes),
		"created_at": time.Now(),
		"roleId":     "665ce8afd343136949deade1",
		"token":      id.String(),
	}
	_, err = mongodb.Post[model.User](newRecord)
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

func PatchPetOwner(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"

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
		http.Error(response, errorMessage[1], http.StatusBadRequest)
		return
	}

	err1 := bcrypt.CompareHashAndPassword([]byte(isSameUser.Password), []byte(requestBody.Password))
	if err1 != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Invalid Password")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var newRecord = bson.M{
		"firstName": requestBody.FirstName,
		"lastName":  requestBody.LastName,
		"email":     requestBody.Email,
		"phoneNo":   requestBody.PhoneNo,
		"password":  requestBody.Password,
	}

	_, err = mongodb.Patch[model.User](filter, newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully updated pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeletePetOwner(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"

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

	_, err := mongodb.Delete[model.User](filter)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully deleted pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
