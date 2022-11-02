package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organizer struct {
	mongoID     primitive.ObjectID `bson:"_id"`
	OrganizerID string             `json:"organizer_id" bson:"organizer_id"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
}

func NewOrganizer(organizerID string, username string, email string) *Organizer {
	return &Organizer{
		OrganizerID: organizerID,
		Username:    username,
		Email:       email,
	}
}
