package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	PhoneNo   string             `json:"phoneNo" bson:"phoneNo"`
	Password  string             `json:"password" bson:"password"`
	RoleId    string             `json:"roleId" bson:"roleId"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Doctor struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    string             `json:"firstName" bson:"firstName"`
	LastName     string             `json:"lastName" bson:"lastName"`
	Email        string             `json:"email" bson:"email"`
	PhoneNo      string             `json:"phoneNo" bson:"phoneNo"`
	Password     string             `json:"password" bson:"password"`
	FatherName   string             `json:"fatherName" bson:"fatherName"`
	Registration string             `json:"registration" bson:"registration"`
	ClinicName   string             `json:"clinicName" bson:"clinicName"`
	Location     Location           `json:"location" bson:"location"`
	Experience   string             `json:"experience" bson:"experience"`
	Bio          string             `json:"bio" bson:"bio"`
	Status       string             `json:"accountStatus" bson:"accountStatus"`
	Token        string             `json:"token" bson:"token"`
	RoleId       string             `json:"roleId" bson:"roleId"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	SenderId   string             `json:"senderId" bson:"senderId"`
	ReceiverId string             `json:"receiverId" bson:"receiverId"`
	Message    string             `json:"message" bson:"message"`
	Rating     int16              `json:"rating" bson:"rating"`
	Token      string             `json:"token" bson:"token"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Token       string             `json:"token" bson:"token"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
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
	Weight     string             `json:"weight" bson:"weight"`
	DietPlan   string             `json:"dietPlan" bson:"dietPlan"`
	Vaccinated bool               `json:"vaccinated" bson:"vaccinated"`
	Tags       []string           `json:"tags" bson:"tags"`
	UserId     string             `json:"userId" bson:"userId"`
	Token      string             `json:"token" bson:"token"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type PetSell struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    string             `json:"userId" bson:"userId"`
	PetId     string             `json:"petId" bson:"petId"`
	Price     float32            `json:"price" bson:"price"`
	ContactNo string             `json:"contactNo" bson:"contactNo"`
	PetDetail *Pets              `json:"petDetail" bson:"petDetail"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
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
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserId     string             `json:"userId" bson:"userId"`
	ReceiverId string             `json:"receiverId" bson:"receiverId"`
	RoomId     string             `json:"roomId" bson:"roomId"`
	Token      string             `json:"token" bson:"token"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
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

type Appointment struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	DoctorID        string             `json:"doctorId" bson:"doctorId"`
	PetOwnerID      string             `json:"petOwnerId" bson:"petOwnerId"`
	AppointmentDate time.Time          `json:"appointmentDate"  bson:"appointmentDate"`
	IsAccepted      bool               `json:"isAccepted" bson:"isAccepted"`
	Status          string             `json:"status" bson:"status"`
	Type            string             `json:"type" bson:"type"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}
