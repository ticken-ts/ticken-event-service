package models

import "github.com/google/uuid"

type Organizer struct {
	OrganizerID uuid.UUID `bson:"organizer_id"`
	Firstname   string    `bson:"firstname"`
	Lastname    string    `bson:"lastname"`
	Username    string    `bson:"username"`
	Email       string    `bson:"email"`
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
