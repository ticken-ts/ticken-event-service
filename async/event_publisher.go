package async

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"ticken-event-service/infra"
	"ticken-event-service/infra/bus"
	"ticken-event-service/models"
)

const (
	NewEventMessageType = "new_event"
)

type eventDTO struct {
	EventID      uuid.UUID `json:"event_id"`
	OrganizerID  uuid.UUID `json:"organizer_id"`
	PvtBCChannel string    `json:"pvt_bc_channel"`
	PubBCAddress string    `json:"pub_bc_address"`
}

type EventPublisher struct {
	busPublisher infra.BusPublisher
}

func NewEventPublisher(busPublisher infra.BusPublisher) *EventPublisher {
	return &EventPublisher{busPublisher: busPublisher}
}

func (processor *EventPublisher) PublishNewEvent(event *models.Event) error {
	dto := mapToDTO(event)

	serializedDTO, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = processor.busPublisher.Publish(
		context.Background(),
		bus.Message{Type: NewEventMessageType, Data: serializedDTO},
	)
	if err != nil {
		return err
	}

	return nil
}

func mapToDTO(event *models.Event) *eventDTO {
	return &eventDTO{
		EventID:      event.EventID,
		OrganizerID:  event.OrganizationID,
		PvtBCChannel: event.PvtBCChannel,
		PubBCAddress: event.PubBCAddress,
	}
}
