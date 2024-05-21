package custom_middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		authorization := request.Header.Get("Authorization")
		if authorization == "" {
			jsonData, err := helping.JsonEncode("Missing Authorization header")
			if err != nil {
				log.Fatal(err)
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("Internal Server Error"))
				return
			}
			response.WriteHeader(http.StatusUnauthorized)
			response.Write(jsonData)
			return
		}
		splitToken := strings.Split(authorization, "Bearer ")
		if len(splitToken) != 2 {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknStr := splitToken[1]

		claims := &data_type.Claims{}
		jwtKey := os.Getenv("JWT_SECRET")

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte("Unauthorized Request"))
				return
			}
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Internal Server Error"))
			return
		}

		if !tkn.Valid {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Unauthorized Request"))
			return
		}

		// expirationString := claims.ExpiresAt
		// expirationTime, err := time.Parse(time.RFC3339, expirationString.String())
		// if err != nil {
		// 	response.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		// if time.Now().UTC().After(expirationTime) {
		// 	response.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		ctx := context.WithValue(request.Context(), "claims", claims)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
