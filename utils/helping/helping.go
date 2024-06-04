package helping

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	data_type "vetner360-backend/utils/type"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func JwtGenerator(response http.ResponseWriter, creds *data_type.Credentials, password string, expirationTime time.Time) (string, error) {
	jwtKey := os.Getenv("JWT_SECRET")

	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password))

	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("password does not match")
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
