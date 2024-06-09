package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	//Common information between each users
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	PhoneNo   string `json:"phoneNo" bson:"phoneNo"`
	Password  string `json:"password" bson:"password"`
	RoleId    string `json:"roleId" bson:"roleId"`
	Token     string `json:"token" bson:"token"`

	//Information for doctor
	FatherName    *string `json:"fatherName" bson:"fatherName"`
	Registration  *string `json:"registration" bson:"registration"`
	ClinicAddress *string `json:"clinicAddress" bson:"clinicAddress"`
	Longitude     *string `json:"longitude" bson:"longitude"`
	Latitude      *string `json:"latitude" bson:"latitude"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	Token       string             `json:"token" bson:"token"`
	// PermissionsIdList *[]primitive.ObjectID `json:"permissionIdList" bson:"permissionIdList"`
}

type Permission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Path        string             `json:"path" bson:"path"`
	Token       string             `json:"token" bson:"token"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type ContactMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Message   string             `json:"message" bson:"message"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Pets struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	NickName   string             `json:"nickName" bson:"nickName"`
	Gender     string             `json:"gender" bson:"gender"`
	BirthDate  time.Time          `json:"birthDate" bson:"birthDate"`
	Type       string             `json:"type" bson:"type"`
	ImagePath  string             `json:"imagePath" bson:"imagePath"`
	Breed      string             `json:"breed" bson:"breed"`
	Note       string             `json:"note" bson:"note"`
	Age        string             `json:"age" bson:"age"`
	Weight     string             `json:"weight" bson:"weight"`
	DietPlan   string             `json:"dietPlan" bson:"dietPlan"`
	Vaccinated bool               `json:"vaccinated" bson:"vaccinated"`
	UserId     string             `json:"userId" bson:"userId"`
	Token      string             `json:"token" bson:"token"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type Activity struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Note      string             `json:"note" bson:"note"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   time.Time          `json:"endTime" bson:"endTime"`
	Status    string             `json:"status" bson:"status"`
	PetId     string             `json:"petId" bson:"petId"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Participant struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    string             `json:"userId" bson:"userId"`
	RoomId    string             `json:"roomId" bson:"roomId"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	SenderId  string             `json:"senderId" bson:"senderId"`
	RoomId    string             `json:"roomId" bson:"roomId"`
	Content   string             `json:"content" bson:"content"`
	Type      string             `json:"type" bson:"type"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
