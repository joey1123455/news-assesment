package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateUser struct {
	Username   string    `json:"user_name" bson:"user omitempty"`
	Prefrences Prefrence `json:"prefrence" bson:"prefrence omitempty"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateUser struct {
	ID         primitive.ObjectID `json:"id" bson:"_id" binding:"required"`
	Username   string             `json:"user_name" bson:"user"`
	Prefrences Prefrence          `json:"prefrence" bson:"prefrence omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserProfile struct {
	ID         primitive.ObjectID `json:"id" bson:"_id" binding:"required"`
	Username   string             `json:"user_name" bson:"user"`
	Prefrences Prefrence          `json:"prefrence" bson:"prefrence"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Prefrence struct {
	Categories []string             `json:"categories" bson:"categories"`
	Liked      []primitive.ObjectID `json:"liked" bson:"liked"`
}

func FilteredResponse(user UserProfile) UserProfile {
	return UserProfile{
		Prefrences: user.Prefrences,
		Username:   user.Username,
	}
}
