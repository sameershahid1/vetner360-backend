package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"
)

func WebSignIn(response http.ResponseWriter, request *http.Request) {
	expirationTime := time.Now().Add(time.Hour * 24 * 7)
	tokenString, err := helping.JwtGenerator(response, request, expirationTime)
	if err != nil {
		log.Fatal(err.Error())
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	}
	http.SetCookie(response, cookie)

	var signData = data_type.SignInType{Message: "Successfully Login"}
	var jsonData, err1 = json.Marshal(signData)

	if err1 != nil {
		log.Fatal(err1.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func MobileSignIn(response http.ResponseWriter, request *http.Request) {
	expirationTime := time.Now().Add(time.Hour * 24 * 7)
	tokenString, err := helping.JwtGenerator(response, request, expirationTime)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var signData = data_type.SignInType{Message: "Successfully Login", Token: &tokenString}
	var jsonData, err1 = json.Marshal(signData)

	if err1 != nil {
		log.Fatal(err1.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func Registration(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Testing testing"))
}
