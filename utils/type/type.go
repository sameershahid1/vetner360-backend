package data_type

import (
	"vetner360-backend/model"

	"github.com/golang-jwt/jwt/v5"
)

type RecordType interface {
	model.User | model.Role | model.Permission | model.ContactMessage | model.Pets | model.Participant | model.Message | model.Activity
}

type PaginationType[T RecordType] struct {
	Page   uint16 `json:"page" validate:"required,number"`
	Limit  uint16 `json:"limit" validate:"required,number"`
	Record *[]T   `json:"record"`
}

type PetOwnerRequestType struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	PhoneNo   string `json:"phoneNo" validate:"required"`
	Password  string `json:"password" validate:"required"`
	UserType  int    `json:"userType" validate:"required"`
}

type PetPostRequestType struct {
	UserId     string `json:"userId" validate:"required"`
	Name       string `json:"name" validate:"required,min=3,max=25"`
	NickName   string `json:"nickName" validate:"required,min=3,max=25"`
	Gender     string `json:"gender" validate:"required"`
	BirthDate  string `json:"birthDate" validate:"required"`
	Image      string `json:"image" validate:"required"`
	Note       string `json:"note" validate:"required,min=3"`
	Age        string `json:"age" validate:"required,number"`
	Weight     string `json:"weight" validate:"required,number"`
	DietPlan   string `json:"dietPlan" validate:"required,min=3"`
	Vaccinated bool   `json:"vaccinated" validate:"omitempty"`
}

type ActivityPostRequestType struct {
	Name      string `json:"name" validate:"required,min=3,max=25"`
	Note      string `json:"note" validate:"required,min=3"`
	StartTime string `json:"startTime" validate:"required"`
	EndTime   string `json:"endTime" validate:"required"`
	PetId     string `json:"petId" validate:"required"`
}

type PetPatchRequestType struct {
	UserId     string `json:"userId" validate:"required"`
	Name       string `json:"name" validate:"required,min=3,max=25"`
	NickName   string `json:"nickName" bson:"nickName"`
	Gender     string `json:"gender" bson:"gender"`
	BirthDate  string `json:"birthDate" bson:"birthDate"`
	Note       string `json:"note" validate:"required,min=3"`
	Age        string `json:"age" validate:"required,number"`
	Weight     string `json:"weight" validate:"required,number"`
	DietPlan   string `json:"dietPlan" validate:"required,min=3"`
	Vaccinated bool   `json:"vaccinated" validate:"omitempty"`
}

type RoleRequestType struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DoctorRequestType struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=25"`
	LastName  string `json:"lastName" validate:"required,min=3,max=25"`
	Email     string `json:"email" validate:"required,email"`
	PhoneNo   string `json:"phoneNo" validate:"required"`
	Password  string `json:"password" validate:"required,min=6,max=25"`

	FatherName    *string `json:"fatherName" validate:"required,min=3,max=25"`
	Registration  *string `json:"registration" validate:"required"`
	ClinicAddress *string `json:"clinicAddress" validate:"required,min=3,max=25"`
	// Longitude     *string `json:"longitude" validate:"required"`
	// Latitude      *string `json:"latitude" validate:"required"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=25"`
}

type Response[T RecordType] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Records *[]T   `json:"records,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

type SignInType struct {
	Message  string  `json:"message"`
	Token    *string `json:"token"`
	UserId   string  `json:"userId"`
	RoleType int     `json:"roleType"`
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
