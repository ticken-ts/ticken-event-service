package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
	"time"
)

func MapEventToEventDTO(event *models.Event) *dto.EventDTO {
	return &dto.EventDTO{
		EventID:  event.EventID.String(),
		Name:     event.Name,
		Date:     event.Date.Format(time.RFC3339),
		OnChain:  event.OnChain,
		Sections: MapSectionListToDTO(event.Sections),
		Poster:   event.PosterAssetID.String(),
	}
}

func MapEventListToDTO(events []*models.Event) []*dto.EventDTO {
	dtos := make([]*dto.EventDTO, len(events))
	for i, e := range events {
		dtos[i] = MapEventToEventDTO(e)
	}
	return dtos
}
