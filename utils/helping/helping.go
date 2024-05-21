package helping

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	static_data "vetner360-backend/utils/data"
	data_type "vetner360-backend/utils/type"

	"github.com/golang-jwt/jwt/v5"
)

func InternalServerError(response http.ResponseWriter, err error) {
	log.Fatal(err)
	response.WriteHeader(http.StatusInternalServerError)
	jsonData, err := JsonEncode("Missing Authorization header")
	if err != nil {
		response.Write([]byte("Internal server error"))
	}
	response.Write(jsonData)
}

func JsonEncode(message string) ([]byte, error) {
	auth := data_type.UnAuthorizeResponse{Message: message}
	jsonData, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func JwtGenerator(response http.ResponseWriter, request *http.Request, expirationTime time.Time) (string, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	var creds data_type.Credentials
	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return "", err
	}

	password, ok := static_data.Users[creds.Email]
	if !ok || password != creds.Password {
		response.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Password does not match")
	}

	claims := &data_type.Claims{
		Username: creds.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return "", err
	}

	return tokenString, nil
}
