package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	mongoID        primitive.ObjectID `bson:"_id"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	Peers          []string           `json:"peers" bson:"peers"`
	Users          []string           `json:"users" bson:"users"`
}

func NewOrganization(id string, peers []string, users []string) *Organization {
	return &Organization{
		OrganizationID: id,
		Peers:          peers,
		Users:          users,
	}
}
