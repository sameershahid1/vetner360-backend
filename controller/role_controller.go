package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoles(response http.ResponseWriter, request *http.Request) {

	var requestBody data_type.PaginationType[model.Role]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var filter = bson.M{}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Role](&filter, &opts, "roles")
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
	id := uuid.New()
	var requestBody data_type.RoleRequestType
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}
	opts := options.FindOneOptions{}
	isSame, _ := mongodb.GetOne[model.Role](bson.M{"name": requestBody.Name}, &opts, "roles")
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
	_, err = mongodb.Post[model.Role](newRecord, "roles")
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
	opts := options.FindOneOptions{}
	isSame, _ := mongodb.GetOne[model.Role](filter, &opts, "roles")
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

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var newRecord = bson.M{
		"name":        requestBody.Name,
		"description": requestBody.Description,
	}

	_, err = mongodb.Patch[model.Role](filter, newRecord, "roles")
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
	opts := options.FindOneOptions{}
	isSame, _ := mongodb.GetOne[model.Role](filter, &opts, "roles")
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

	_, err := mongodb.Delete[model.Role](filter, "roles")
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
