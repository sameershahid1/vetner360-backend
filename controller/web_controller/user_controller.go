package web_controller

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
	"golang.org/x/crypto/bcrypt"
)

func GetUser(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	if query["userType"] == nil {
		helping.InternalServerError(response, errors.New("missing userType query"), http.StatusBadRequest)
		return
	}

	if query["userType"][0] != "petOwner" && query["userType"][0] != "guest" {
		helping.InternalServerError(response, errors.New("incorrect userType query"), http.StatusBadRequest)
		return
	}

	var requestBody data_type.PaginationType[model.User]
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

	userType := query["userType"][0]
	roleID := ""
	if userType == "petOwner" {
		roleID = "665ceb8baf682359fe5990a8"
	}
	if userType == "guest" {
		roleID = "665cecbdc6206b06eddaaccb"
	}

	var filter = bson.M{"roleId": roleID}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.User](&filter, &opts, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	total, err := mongodb.TotalDocs[model.User](&filter, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request", Records: &records, Count: &total}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PostUser(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
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

	roleID := ""
	if requestBody.UserType == 1 {
		roleID = "665ceb8baf682359fe5990a8"
	} else if requestBody.UserType == 2 {
		roleID = "665cecbdc6206b06eddaaccb"
	} else {
		helping.InternalServerError(response, errors.New("invalid userId"), http.StatusBadRequest)
		return
	}

	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.User](bson.M{"email": requestBody.Email}, &opts, "users")
	if isSameUser != nil {
		helping.InternalServerError(response, errors.New("user already exists"), http.StatusBadRequest)
		return
	}

	cost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), cost)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var newRecord = bson.M{
		"firstName":  requestBody.FirstName,
		"lastName":   requestBody.LastName,
		"email":      requestBody.Email,
		"phoneNo":    requestBody.PhoneNo,
		"password":   string(bytes),
		"created_at": time.Now(),
		"roleId":     roleID,
		"token":      id.String(),
	}
	_, err = mongodb.Post[model.User](newRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Completed Request"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func PatchUser(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}

	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.User](filter, &opts, "users")
	if isSameUser == nil {
		helping.InternalServerError(response, errors.New("user does not exists"), http.StatusBadRequest)
		return
	}

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

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully updated pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	var id = chi.URLParam(request, "id")
	var filter = bson.M{"token": id}
	opts := options.FindOneOptions{}
	isSameUser, _ := mongodb.GetOne[model.User](filter, &opts, "users")
	if isSameUser == nil {
		helping.InternalServerError(response, errors.New("user does not exists"), http.StatusBadRequest)
		return
	}

	_, err := mongodb.Delete[model.User](filter, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully deleted pet owner"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
