package web_controller

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

func GetContactMessage(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.ContactMessage]
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

	records, err := mongodb.GetAll[model.ContactMessage](&filter, &opts, "contactmessages")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.ContactMessage]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostContactMessage(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.ContactMessageType
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
		"email":      requestBody.Email,
		"message":    requestBody.Message,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.ContactMessage](newRecord, "contactmessages")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.ContactMessage]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchContactMessage(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}

	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.ContactMessage](filter, &opts, "contactmessages")
	if isSameUser == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Contact Message does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var requestBody data_type.ContactMessageType
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

	var updateRecord = bson.M{
		"email":   requestBody.Email,
		"message": requestBody.Message,
	}

	_, err = mongodb.Patch[model.User](filter, updateRecord, "contactmessages")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully updated contact message"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteContactMessage(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.ContactMessage](filter, &opts, "contactmessages")
	if isSameUser == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Contact Message does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	_, err := mongodb.Delete[model.ContactMessage](filter, "contactmessages")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully deleted contact message"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
