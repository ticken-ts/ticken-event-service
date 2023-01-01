package models

import "github.com/google/uuid"

type Organizer struct {
	OrganizerID uuid.UUID `json:"organizer_id" bson:"organizer_id"`
	Firstname   string    `json:"firstname" bson:"firstname"`
	Lastname    string    `json:"lastname" bson:"lastname"`
	Username    string    `json:"username" bson:"username"`
	Email       string    `json:"email" bson:"email"`
}

func NewOrganizer(organizerID uuid.UUID, firstname, lastname, username, email string) *Organizer {
	return &Organizer{
		OrganizerID: organizerID,
		Firstname:   firstname,
		Lastname:    lastname,
		Username:    username,
		Email:       email,
	}
}
