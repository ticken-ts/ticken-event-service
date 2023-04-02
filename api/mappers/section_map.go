package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
)

func MapSectionToDTO(section *models.Section) *dto.SectionDTO {
	return &dto.SectionDTO{
		EventID:      section.EventID.String(),
		Name:         section.Name,
		TotalTickets: section.TotalTickets,
		Price:        section.TicketPrice,
	}
}

func MapSectionListToDTO(sections []*models.Section) []*dto.SectionDTO {
	dtos := make([]*dto.SectionDTO, len(sections))
	for i, section := range sections {
		dtos[i] = MapSectionToDTO(section)
	}
	return dtos
}
