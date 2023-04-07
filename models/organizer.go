package models

import "github.com/google/uuid"

type Organizer struct {
	OrganizerID uuid.UUID `bson:"organizer_id"`
	Firstname   string    `bson:"firstname"`
	Lastname    string    `bson:"lastname"`
	Username    string    `bson:"username"`
	Email       string    `bson:"email"`
}
