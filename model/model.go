package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	//Common information between each users
	FirstName string  `json:"firstName" bson:"firstName"`
	LastName  string  `json:"lastName" bson:"lastName"`
	Email     string  `json:"email" bson:"email"`
	PhoneNo   string  `json:"phoneNo" bson:"phoneNo"`
	Password  string  `json:"password" bson:"password"`
	RoleId    *string `json:"roleId" bson:"roleId"`
	Token     string  `json:"token" bson:"token"`

	//Pet-owner User Pet list
	PetIdList *[]string `json:"petIdList" bson:"petIdList"`

	//Information for doctor
	FatherName    *string `json:"fatherName" bson:"fatherName"`
	Registration  *string `json:"registration" bson:"registration"`
	ClinicAddress *string `json:"clinicAddress" bson:"clinicAddress"`
	Longitude     *string `json:"longitude" bson:"longitude"`
	Latitude      *string `json:"latitude" bson:"latitude"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Role struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	Name              string             `json:"name" bson:"name"`
	Description       string             `json:"description" bson:"description"`
	PermissionsIdList *[]string          `json:"permissionIdList" bson:"permissionIdList"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

type Permission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Path        string             `json:"path" bson:"path"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ContactMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Pets struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Type      string             `json:"type" bson:"type"`
	Bread     string             `json:"bread" bson:"bread"`
	Age       string             `json:"age" bson:"age"`
	Weight    string             `json:"weight" bson:"weight"`
	DietPlan  string             `json:"dietPlan" bson:"dietPlan"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Activity struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   time.Time          `json:"endTime" bson:"endTime"`
	Status    string             `json:"status" bson:"status"`
	PetId     *string            `json:"petId" bson:"petId"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Chat struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	SenderId   string             `json:"senderId" bson:"senderId"`
	ReceiverId string             `json:"receiverId" bson:"receiverId"`
	Message    string             `json:"name" bson:"message"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
