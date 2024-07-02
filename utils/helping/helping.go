package helping

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	data_type "vetner360-backend/utils/type"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func InternalServerError(response http.ResponseWriter, err error, statusCode int) {
	response.WriteHeader(statusCode)
	jsonData, err := JsonEncode(err.Error())
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

func GetValidator() *validator.Validate {
	return validator.New()
}

func ValidatingData(data interface{}, response http.ResponseWriter, validate *validator.Validate) error {
	err := validate.Struct(data)
	if err != nil {
		errorMessageList := strings.Split(err.Error(), "\n")
		errorMessage := strings.Split(errorMessageList[0], "Error:")
		response.WriteHeader(http.StatusBadRequest)
		jsonErrorMessage, err := JsonEncode(errorMessage[1])
		if err != nil {
			response.Write([]byte("Internal server side error"))
		}
		response.Write(jsonErrorMessage)
		return err
	}

	return nil
}
