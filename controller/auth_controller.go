package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(response http.ResponseWriter, request *http.Request) {
	var creds data_type.Credentials
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"
	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		jsonMessage, err := helping.JsonEncode("Internal server error")
		if err != nil {
			response.Write([]byte("Internal server error"))
			return
		}
		response.Write(jsonMessage)
		return
	}

	var filter = bson.M{"email": creds.Email}
	record, err := mongodb.GetOne[model.User](filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		jsonMessage, err := helping.JsonEncode("Invalid email")
		if err != nil {
			response.Write([]byte("Internal server error"))
			return
		}
		response.Write(jsonMessage)
		return
	}

	if record == nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonMessage, err := helping.JsonEncode("Invalid email")
		if err != nil {
			response.Write([]byte("Internal server error"))
			return
		}
		response.Write(jsonMessage)
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24 * 7)
	tokenString, err := helping.JwtGenerator(response, &creds, record.Password, expirationTime)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		jsonMessage, err := helping.JsonEncode(err.Error())
		if err != nil {
			response.Write([]byte("Internal server error"))
			return
		}
		response.Write(jsonMessage)
		return
	}

	var signData = data_type.SignInType{Message: "Successfully Login", Token: &tokenString, UserId: record.Token}
	var jsonData, err1 = json.Marshal(signData)

	if err1 != nil {
		log.Fatal(err1.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func PetOwnerORGuestRegistration(response http.ResponseWriter, request *http.Request) {
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
		jsonResponse, err := helping.JsonEncode("Account already exists")
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
		"token":      id.String(),
		"roleId":     "665ceb8baf682359fe5990a8",
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.User](newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully registered User"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func DoctorRegistration(response http.ResponseWriter, request *http.Request) {
	mongodb.Database = "vetner360"
	mongodb.Collection = "users"
	id := uuid.New()
	var requestBody data_type.DoctorRequestType
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
		jsonResponse, err := helping.JsonEncode("Doctor already exists")
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
		"firstName":     requestBody.FirstName,
		"lastName":      requestBody.LastName,
		"email":         requestBody.Email,
		"phoneNo":       requestBody.PhoneNo,
		"password":      string(bytes),
		"fatherName":    requestBody.FirstName,
		"registration":  requestBody.Registration,
		"clinicAddress": requestBody.ClinicAddress,
		"token":         id.String(),
		"roleId":        "665cec7fc6206b06eddaacca",
		"created_at":    time.Now(),
	}
	_, err = mongodb.Post[model.User](newRecord)
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.User]{Status: true, Message: "Successfully Register Doctor"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
