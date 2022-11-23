package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	mongoID primitive.ObjectID `bson:"_id"`
}

func NewOrganization() *Organization {
	return &Organization{}
}
