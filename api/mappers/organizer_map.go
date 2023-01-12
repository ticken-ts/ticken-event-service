package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
)

func MapOrganizerToOrganizerDTO(organizer *models.Organizer) *dto.OrganizerDTO {
	return &dto.OrganizerDTO{
		OrganizerID: organizer.OrganizerID,
		Firstname:   organizer.Firstname,
		Lastname:    organizer.Lastname,
		Username:    organizer.Username,
		Email:       organizer.Email,
	}
}
