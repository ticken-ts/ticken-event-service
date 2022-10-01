package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	mongoID        primitive.ObjectID `bson:"_id"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	Events         []string           `json:"events" bson:"events"`
	Peers          []string           `json:"peers" bson:"peers"`
	Users          []string           `json:"users" bson:"users"`
}

func NewOrganization(id string, events []string, peers []string, users []string) *Organization {
	return &Organization{
		OrganizationID: uuid.NewString(),
		Events:         events,
		Peers:          peers,
		Users:          users,
	}
}
