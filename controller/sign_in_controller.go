package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var users = map[string]string{
	"sameer": "123456",
	"user2":  "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Signin(response http.ResponseWriter, request *http.Request) {
	jwtKey := os.Getenv("JWT_SECRET")

	var creds Credentials
	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	password, ok := users[creds.Username]
	if !ok || password != creds.Password {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(100 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}

	http.SetCookie(response, cookie)
}
