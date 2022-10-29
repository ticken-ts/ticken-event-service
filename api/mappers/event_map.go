package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
	"time"
)

func MapEventToCreatedTicketDTO(event *models.Event) *dto.CreatedEventDTO {
	return &dto.CreatedEventDTO{
		EventID: event.EventID,
		Name:    event.Name,
		Date:    event.Date.Format(time.RFC3339),
		OnChain: event.OnChain,
	}
}
