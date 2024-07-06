package mobile_controller

import (
	"encoding/json"
	"errors"
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

func GetMySellPets(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.PetSell]
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

	var userId = chi.URLParam(request, "userId")
	var filter = bson.M{
		"userId": userId,
	}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.PetSell](&filter, &opts, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}
	if records == nil {
		records = []model.PetSell{}
	}

	petIdList := []string{}
	for x := range records {
		petIdList = append(petIdList, records[x].PetId)
	}

	filter = bson.M{
		"userId": userId,
		"token":  bson.M{"$in": petIdList},
	}
	petRecords, err := mongodb.GetAll[model.Pets](&filter, &opts, "pets")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	if petRecords == nil {
		petRecords = []model.Pets{}
	}

	total, err := mongodb.TotalDocs[model.PetSell](&bson.M{"userId": userId}, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Pets]{Status: true, Message: "Successfully Completed Request", Records: &petRecords, Count: &total}
	jsonData, err := json.Marshal(requestResponse)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostMyPetToSell(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.PetSellRequestType
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

	var filter = bson.M{"userId": requestBody.UserId, "petId": requestBody.PetId}
	opts := options.FindOneOptions{}
	isSamePets, _ := mongodb.GetOne[model.PetSell](filter, &opts, "petsells")
	if isSamePets != nil {
		helping.InternalServerError(response, errors.New("pet is already on sell"), http.StatusBadRequest)
		return
	}

	var newRecord = bson.M{
		"userId":     requestBody.UserId,
		"petId":      requestBody.PetId,
		"price":      requestBody.Price,
		"contactNo":  requestBody.ContactNo,
		"token":      id.String(),
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.PetSell](newRecord, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.PetSell]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchMyPetOnSell(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	var requestBody data_type.PetSellRequestType
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

	var newRecord = bson.M{
		"price":     requestBody.Price,
		"contactNo": requestBody.ContactNo,
	}
	_, err = mongodb.Patch[model.PetSell](filter, newRecord, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.PetSell]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteMyPetOnSell(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSamePets, _ := mongodb.GetOne[model.PetSell](filter, &opts, "petsells")
	if isSamePets == nil {
		helping.InternalServerError(response, errors.New("pet does not exists"), http.StatusInternalServerError)
		return
	}

	_, err := mongodb.Delete[model.PetSell](filter, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.PetSell]{Status: true, Message: "Successfully deleted pet on sell"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
