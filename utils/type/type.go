package data_type

import (
	"vetner360-backend/model"

	"github.com/golang-jwt/jwt/v5"
)

type RecordType interface {
	model.User | model.Role | model.Permission | model.ContactMessage | model.Pets
}

type PaginationType[T RecordType] struct {
	Page   uint16 `json:"page" validate:"required"`
	Limit  uint16 `json:"limit" validate:"required"`
	Record *[]T   `json:"record"`
}

type PetOwnerRequestType struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	PhoneNo   string `json:"phoneNo" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response[T RecordType] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Records *[]T   `json:"records,omitempty"`
}

type SignInType struct {
	Message string  `json:"message"`
	Token   *string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UnAuthorizeResponse struct {
	Message string `json:"message"`
}
