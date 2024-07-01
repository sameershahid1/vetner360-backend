package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(response http.ResponseWriter, request *http.Request) {
	var creds data_type.Credentials
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

	validate := helping.GetValidator()
	err = helping.ValidatingData(creds, response, validate)
	if err != nil {
		return
	}

	var filter = bson.M{"email": creds.Email}
	opts := options.FindOneOptions{}
	record, err := mongodb.GetOne[model.User](filter, &opts, "users")
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

	var roleType int = 0
	if record.RoleId == "665ceb8baf682359fe5990a8" {
		roleType = 1
	} else if record.RoleId == "665cecbdc6206b06eddaaccb" {
		roleType = 2
	} else if record.RoleId == "665cec7fc6206b06eddaacca" {
		roleType = 3
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

	var signData = data_type.SignInType{Message: "Successfully Login", Token: &tokenString, UserId: record.Token, RoleType: roleType}
	var jsonData, err1 = json.Marshal(signData)

	if err1 != nil {
		log.Fatal(err1.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func UserRegistration(response http.ResponseWriter, request *http.Request) {
	id := uuid.New()
	var requestBody data_type.PetOwnerRequestType
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
	isSameUser, _ := mongodb.GetOne[model.User](bson.M{"email": requestBody.Email}, &opts, "users")
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

	roleID := ""
	if requestBody.UserType == 1 {
		roleID = "665ceb8baf682359fe5990a8"
	} else {
		roleID = "665cecbdc6206b06eddaaccb"
	}

	var newRecord = bson.M{
		"firstName":  requestBody.FirstName,
		"lastName":   requestBody.LastName,
		"email":      requestBody.Email,
		"phoneNo":    requestBody.PhoneNo,
		"password":   string(bytes),
		"token":      id.String(),
		"roleId":     roleID,
		"created_at": time.Now(),
	}
	_, err = mongodb.Post[model.User](newRecord, "users")
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
	id := uuid.New()
	var requestBody data_type.DoctorRequestType
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
	isSameUser, _ := mongodb.GetOne[model.Doctor](bson.M{"email": requestBody.Email}, &opts, "users")
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

	var location = model.Location{
		Type:        "Point",
		Coordinates: []float64{requestBody.Longitude, requestBody.Latitude},
	}

	var newRecord = bson.M{
		"firstName":    requestBody.FirstName,
		"lastName":     requestBody.LastName,
		"email":        requestBody.Email,
		"phoneNo":      requestBody.PhoneNo,
		"password":     string(bytes),
		"fatherName":   requestBody.FatherName,
		"registration": requestBody.Registration,
		"clinicName":   requestBody.ClinicName,
		"experience":   requestBody.Experience,
		"bio":          requestBody.Bio,
		"location":     location,
		"token":        id.String(),
		"roleId":       "665cec7fc6206b06eddaacca",
		"created_at":   time.Now(),
	}
	_, err = mongodb.Post[model.Doctor](newRecord, "users")
	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully Register Doctor"}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
