package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
)

func MapSectionToDTO(section *models.Section) *dto.SectionDTO {
	return &dto.SectionDTO{
		EventID:      section.EventID,
		Name:         section.Name,
		TotalTickets: section.TotalTickets,
		OnChain:      section.OnChain,
	}
}
