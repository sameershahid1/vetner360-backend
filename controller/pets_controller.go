package controller

import (
	"encoding/json"
	"fmt"
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

func GetMyPetList(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "pets"

	var requestBody data_type.PaginationType[model.Pets]
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

	var id = chi.URLParam(request, "id")
	var filter = bson.M{"userId": id}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Pets](&filter, &opts)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	if records == nil {
		records = []model.Pets{}
	}

	fmt.Println(records)

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetPetDetail(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "pets"

	var requestBody data_type.PaginationType[model.Pets]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var id = chi.URLParam(request, "id")
	var userId = chi.URLParam(request, "userId")
	var filter = bson.M{"userId": userId, "token": id}

	records, err := mongodb.GetOne[model.Pets](filter)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request", Data: records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostPet(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "pets"
	id := uuid.New()
	var requestBody data_type.PetPostRequestType
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

	isSamePets, _ := mongodb.GetOne[model.Pets](bson.M{"name": requestBody.Name})
	if isSamePets != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Pets already exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var newRecord = bson.M{
		"name":       requestBody.Name,
		"type":       "Cat",
		"breed":      "Pakistani",
		"note":       requestBody.Note,
		"age":        requestBody.Age,
		"weight":     requestBody.Weight,
		"dietPlan":   requestBody.DietPlan,
		"vaccinated": requestBody.Vaccinated,
		"userId":     requestBody.UserId,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.Pets](newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchPet(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "pets"

	isSamePets, _ := mongodb.GetOne[model.Pets](filter)
	if isSamePets == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Pets does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var requestBody data_type.PetPatchRequestType
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
		"name": requestBody.Name,
		// "type":       requestBody.Type,
		// "breed":      requestBody.breed,
		"note":       requestBody.Note,
		"age":        requestBody.Age,
		"weight":     requestBody.Weight,
		"dietPlan":   requestBody.DietPlan,
		"vaccinated": requestBody.Vaccinated,
	}

	_, err = mongodb.Patch[model.Pets](filter, newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully updated pet"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeletePet(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	mongodb.Database = "vetner360"
	mongodb.Collection = "pets"

	isSamePets, _ := mongodb.GetOne[model.Pets](filter)
	if isSamePets == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Pets does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	_, err := mongodb.Delete[model.Pets](filter)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully deleted pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
