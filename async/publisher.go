package async

import (
	"fmt"
	"ticken-event-service/infra"
)

type Publisher struct {
	busPublisher infra.BusPublisher
	*EventPublisher
}

func NewPublisher(busPublisher infra.BusPublisher) (*Publisher, error) {
	if !busPublisher.IsConnected() {
		return nil, fmt.Errorf("bus publisher is not connected")
	}

	publisher := &Publisher{
		busPublisher:   busPublisher,
		EventPublisher: NewEventPublisher(busPublisher),
	}

	return publisher, nil
}
