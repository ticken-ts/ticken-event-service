package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
)

func MapValidatorToValidatorDTO(validator *models.Validator) *dto.Validator {
	return &dto.Validator{
		ValidatorID:    validator.ValidatorID.String(),
		CreatedBy:      validator.CreatedBy.String(),
		OrganizationID: validator.OrganizationID.String(),
		Email:          validator.Email,
	}
}
