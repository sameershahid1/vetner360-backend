package data_type

import (
	"vetner360-backend/model"

	"github.com/golang-jwt/jwt/v5"
)

type RecordType interface {
	model.User | model.Doctor | model.Role | model.Permission | model.ContactMessage | model.Pets
}

type Response[T RecordType] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Records []T    `json:"records,omitempty"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type SignInType struct {
	Message string  `json:"message"`
	Token   *string `json:"token"`
}
