package data_type

import (
	"vetner360-backend/model"

	"github.com/golang-jwt/jwt/v5"
)

type RecordType interface {
	model.User | model.Role | model.Permission | model.ContactMessage | model.Pets | model.Participant | model.Message
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

type PetPostRequestType struct {
	UserId     string `json:"userId" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Note       string `json:"note" validate:"required"`
	Age        string `json:"age" validate:"required"`
	Weight     string `json:"weight" validate:"required"`
	DietPlan   string `json:"dietPlan" validate:"required"`
	Vaccinated bool   `json:"vaccinated" validate:"required"`
}

type PetPatchRequestType struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Bread      string `json:"bread" validate:"required"`
	Note       string `json:"note" validate:"required"`
	Age        string `json:"age" validate:"required"`
	Weight     string `json:"weight" validate:"required"`
	DietPlan   string `json:"dietPlan" validate:"required"`
	Vaccinated bool   `json:"vaccinated" validate:"required"`
}

type RoleRequestType struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DoctorRequestType struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	PhoneNo   string `json:"phoneNo" validate:"required"`
	Password  string `json:"password" validate:"required"`

	FatherName    *string `json:"fatherName" validate:"required"`
	Registration  *string `json:"registration" validate:"required"`
	ClinicAddress *string `json:"clinicAddress" validate:"required"`
	// Longitude     *string `json:"longitude" validate:"required"`
	// Latitude      *string `json:"latitude" validate:"required"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response[T RecordType] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Records *[]T   `json:"records,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

type SignInType struct {
	Message string  `json:"message"`
	Token   *string `json:"token"`
	UserId  string  `json:"userId"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UnAuthorizeResponse struct {
	Message string `json:"message"`
}

type EventMessageType struct {
	RoomId     string `json:"roomId" validate:"required"`
	SenderId   string `json:"senderId" validate:"required"`
	ReceiverId string `json:"receiverId" validate:"required"`
	Message    string `json:"message" validate:"required"`
}
