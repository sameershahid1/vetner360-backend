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
)

func GetRoles(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "roles"

	var requestBody data_type.PaginationType[model.Role]
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

	records, err := mongodb.GetAll[model.Role](&filter, &opts)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Role]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostRoleOwner(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "roles"
	id := uuid.New()
	var requestBody data_type.RoleRequestType
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

	isSame, _ := mongodb.GetOne[model.Role](bson.M{"name": requestBody.Name})
	if isSame != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Role already exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var newRecord = bson.M{
		"name":        requestBody.Name,
		"description": requestBody.Description,
		"token":       id,
		"created_at":  time.Now(),
	}
	_, err = mongodb.Post[model.Role](newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Role]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchRoleOwner(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "roles"

	isSame, _ := mongodb.GetOne[model.Role](filter)
	if isSame == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Role does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var requestBody data_type.RoleRequestType
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

	var newRecord = bson.M{
		"name":        requestBody.Name,
		"description": requestBody.Description,
	}

	_, err = mongodb.Patch[model.Role](filter, newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Role]{Status: true, Message: "Successfully updated Role"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteRoleOwner(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "roles"

	isSame, _ := mongodb.GetOne[model.Role](filter)
	if isSame == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Role does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	_, err := mongodb.Delete[model.Role](filter)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Role]{Status: true, Message: "Successfully deleted Role"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
