package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
	"time"
)

func MapEventToCreatedEventDTO(event *models.Event) *dto.EventDTO {
	return &dto.EventDTO{
		EventID: event.EventID,
		Name:    event.Name,
		Date:    event.Date.Format(time.RFC3339),
		OnChain: event.OnChain,
	}
}
