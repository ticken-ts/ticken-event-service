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
	NewEventMessageType          = "new_event"
	UpdateEventStatusMessageType = "update_status"
)

type createEventDTO struct {
	EventID        uuid.UUID `json:"event_id"`
	OrganizerID    uuid.UUID `json:"organizer_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	PvtBCChannel   string    `json:"pvt_bc_channel"`
	PubBCAddress   string    `json:"pub_bc_address"`
}

type updateEventStatusDTO struct {
	EventID uuid.UUID `json:"event_id"`
	Status  string    `json:"status"`
}

type EventPublisher struct {
	busPublisher infra.BusPublisher
}

func NewEventPublisher(busPublisher infra.BusPublisher) *EventPublisher {
	return &EventPublisher{busPublisher: busPublisher}
}

func (processor *EventPublisher) PublishNewEvent(event *models.Event) error {
	dto := &createEventDTO{
		EventID:        event.EventID,
		OrganizerID:    event.OrganizerID,
		OrganizationID: event.OrganizationID,
		PvtBCChannel:   event.PvtBCChannel,
		PubBCAddress:   event.PubBCAddress,
	}

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

func (processor *EventPublisher) PublishStatusUpdate(event *models.Event) error {
	dto := &updateEventStatusDTO{
		EventID: event.EventID,
		Status:  string(event.Status),
	}

	serializedDTO, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = processor.busPublisher.Publish(
		context.Background(),
		bus.Message{Type: UpdateEventStatusMessageType, Data: serializedDTO},
	)
	if err != nil {
		return err
	}

	return nil
}
