package models

import "github.com/google/uuid"

type Validator struct {
	ValidatorID    uuid.UUID `bson:"validator_id"`
	CreatedBy      uuid.UUID `bson:"created_by"`
	OrganizationID uuid.UUID `bson:"organization_id"`
	Email          string    `bson:"email"`
}
