package middleware

import (
	"context"
	"net/http"
	"os"
	"vetner360-backend/controller"

	"github.com/golang-jwt/jwt/v5"
)

type UnAuthroizeResponse struct {
	Message string `json:"message"`
}

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte("Unauthorize Access"))
				return
			}
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := cookie.Value
		claims := &controller.Claims{}
		jwtKey := os.Getenv("JWT_SECRET")

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
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
			response.Write([]byte("Unauthorize Access"))
			return
		}

		ctx := context.WithValue(request.Context(), "claims", claims)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
