package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	PhoneNo   string             `json:"phoneNo"`
	Password  string             `json:"password"`
	RoleId    *string            `json:"roleId"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Doctor struct {
	User
	FatherName   string `json:"fatherName"`
	Registration string `json:"registration"`
	Longitude    string `json:"longitude"`
	Latitude     string `json:"latitude"`
}

type Role struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	PermissionsID *[]string          `json:"permission"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Permission struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Path        string             `json:"path"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type ContactMessage struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string             `json:"email"`
	Message   string             `json:"message"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Pets struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `json:"name"`
	PetOwnerID primitive.ObjectID `json:"petOwnerId" bson:""`
	Email      string             `json:"email"`
	Message    string             `json:"message"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
