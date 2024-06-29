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

func GetActivityList(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.Activity]
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

	var id = chi.URLParam(request, "petId")
	var filter = bson.M{"petId": id}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Activity](&filter, &opts, "activity")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	if records == nil {
		records = []model.Activity{}
	}

	var requestResponse = data_type.Response[model.Activity]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostActivity(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.ActivityPostRequestType
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

	opts := options.FindOneOptions{}
	isSameActivity, _ := mongodb.GetOne[model.Activity](bson.M{"name": requestBody.Name}, &opts, "activity")
	if isSameActivity != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Activity already exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	layout := "2006-01-02T15:04:05.999"
	parseStartTime, err := time.Parse(layout, requestBody.StartTime)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Error parsing Start date string:" + err.Error())
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}
	parseEndTime, err := time.Parse(layout, requestBody.EndTime)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Error parsing End date string:" + err.Error())
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var newRecord = bson.M{
		"name":       requestBody.Name,
		"note":       requestBody.Note,
		"startTime":  parseStartTime,
		"endTime":    parseEndTime,
		"petId":      requestBody.PetId,
		"status":     "Active",
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.Activity](newRecord, "activity")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Activity]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchActivity(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")

	var requestBody data_type.ActivityPostRequestType
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

	var filter = bson.M{"token": id, "petId": requestBody.PetId}
	opts := options.FindOneOptions{}
	isSameActivity, _ := mongodb.GetOne[model.Activity](filter, &opts, "activity")
	if isSameActivity == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Activity does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	layout := "2006-01-02T15:04:05.999"
	parseStartTime, err := time.Parse(layout, requestBody.StartTime)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Error parsing Start date string:" + err.Error())
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}
	parseEndTime, err := time.Parse(layout, requestBody.EndTime)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Error parsing End date string:" + err.Error())
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	var updateRecord = bson.M{
		"name":      requestBody.Name,
		"note":      requestBody.Note,
		"startTime": parseStartTime,
		"endTime":   parseEndTime,
		"status":    "Active",
	}

	_, err = mongodb.Patch[model.Activity](filter, updateRecord, "activity")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Activity]{Status: true, Message: "Successfully updated activity"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteActivity(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var petId = chi.URLParam(request, "petId")
	var filter = bson.M{"token": id, "petId": petId}
	opts := options.FindOneOptions{}
	isSameActivity, _ := mongodb.GetOne[model.Activity](filter, &opts, "activity")
	if isSameActivity == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := helping.JsonEncode("Activity does not exists")
		if err != nil {
			helping.InternalServerError(response, err)
			return
		}
		response.Write(jsonResponse)
		return
	}

	_, err := mongodb.Delete[model.Activity](filter, "activity")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Activity]{Status: true, Message: "Successfully deleted activity"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
