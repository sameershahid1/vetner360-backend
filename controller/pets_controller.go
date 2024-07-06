package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func GetMyPetList(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.Pets]
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

	var id = chi.URLParam(request, "userId")
	var filter = bson.M{"userId": id}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Pets](&filter, &opts, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	if records == nil {
		records = []model.Pets{}
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetPetDetail(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.Pets]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var id = chi.URLParam(request, "id")
	var userId = chi.URLParam(request, "userId")
	var filter = bson.M{"userId": userId, "token": id}
	opts := options.FindOneOptions{}
	records, err := mongodb.GetOne[model.Pets](filter, &opts, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request", Data: records}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostPet(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.PetPostRequestType
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

	imageBytes, err := base64.StdEncoding.DecodeString(requestBody.Image)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("image-%d.jpg", time.Now().UnixNano())
	filePath := fmt.Sprintf("public/%s", filename)
	err = os.WriteFile(filePath, imageBytes, 0644)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	layout := "2006-01-02T15:04:05.999"
	parsedTime, err := time.Parse(layout, requestBody.BirthDate)
	if err != nil {
		fmt.Println("Error parsing date string:", err)
		return
	}

	var newRecord = bson.M{
		"name":       requestBody.Name,
		"nickName":   requestBody.NickName,
		"gender":     requestBody.Gender,
		"birthDate":  parsedTime,
		"imagePath":  filename,
		"note":       requestBody.Note,
		"weight":     requestBody.Weight,
		"dietPlan":   requestBody.DietPlan,
		"vaccinated": requestBody.Vaccinated,
		"type":       "Cat",
		"breed":      "Persian",
		"tags":       requestBody.Tags,
		"userId":     requestBody.UserId,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.Pets](newRecord, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchPet(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")

	var requestBody data_type.PetPatchRequestType
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
	var filter = bson.M{"token": id, "userId": requestBody.UserId}
	isSamePets, _ := mongodb.GetOne[model.Pets](filter, &opts, "pets")
	if isSamePets == nil {
		helping.InternalServerError(response, errors.New("pet does not exists"), http.StatusInternalServerError)
		return
	}

	var updateRecord = bson.M{
		"name":       requestBody.Name,
		"note":       requestBody.Note,
		"nickName":   requestBody.NickName,
		"gender":     requestBody.Gender,
		"birthDate":  requestBody.BirthDate,
		"age":        requestBody.Age,
		"weight":     requestBody.Weight,
		"dietPlan":   requestBody.DietPlan,
		"vaccinated": requestBody.Vaccinated,
	}

	_, err = mongodb.Patch[model.Pets](filter, updateRecord, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully updated pet"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeletePet(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var userId = chi.URLParam(request, "userId")
	var filter = bson.M{"token": id, "userId": userId}
	opts := options.FindOneOptions{}
	isSamePets, _ := mongodb.GetOne[model.Pets](filter, &opts, "pets")
	if isSamePets == nil {
		helping.InternalServerError(response, errors.New("pet does not exists"), http.StatusInternalServerError)
		return
	}

	_, err := mongodb.Delete[model.Pets](filter, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully deleted pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
