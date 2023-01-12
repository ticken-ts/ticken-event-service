package models

import "github.com/google/uuid"

type Organizer struct {
	OrganizerID string `bson:"organizer_id"`
	Firstname   string `bson:"firstname"`
	Lastname    string `bson:"lastname"`
	Username    string `bson:"username"`
	Email       string `bson:"email"`
}

func NewOrganizer(organizerID uuid.UUID, firstname, lastname, username, email string) *Organizer {
	return &Organizer{
		OrganizerID: organizerID.String(),
		Firstname:   firstname,
		Lastname:    lastname,
		Username:    username,
		Email:       email,
	}
}
