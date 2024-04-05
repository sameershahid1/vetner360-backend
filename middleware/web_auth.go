package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	data_type "vetner360-backend/utils/type"

	"github.com/golang-jwt/jwt/v5"
)

type UnAuthorizeResponse struct {
	Message string `json:"message"`
}

func jsonEncode(message string) ([]byte, error) {
	auth := UnAuthorizeResponse{Message: message}
	jsonData, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func WebSignInVerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("jwt")

		if err != nil {
			if err == http.ErrNoCookie {
				jsonData, err := jsonEncode("Unauthorize Access")
				if err != nil {
					log.Fatal(err.Error())
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				response.WriteHeader(http.StatusUnauthorized)
				response.Write(jsonData)
				return
			}
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := cookie.Value

		claims := &data_type.Claims{}
		jwtKey := os.Getenv("JWT_SECRET")

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "claims", claims)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
