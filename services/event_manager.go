package services

import (
	"fmt"
	"ticken-event-service/blockchain/pvtbc"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type eventManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   pvtbc.TickenConnector
}

func NewEventManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector pvtbc.TickenConnector,
) EventManager {
	newEventMan := new(eventManager)
	newEventMan.eventRepository = eventRepository
	newEventMan.ticketRepository = ticketRepository
	newEventMan.pvtbcConnector = pvtbcConnector
	return newEventMan
}

func (eventManager *eventManager) AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error) {
	event := models.NewEvent(EventID, OrganizerID, PvtBCChannel)
	err := eventManager.eventRepository.AddEvent(event)
	if err != nil {
		return nil, err
	}
	return event, err
}

func (eventManager *eventManager) GetEvent(eventId string, userId string) (*models.Event, error) {
	println("getting event with id:", eventId)

	event := eventManager.eventRepository.FindEvent(eventId)
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}
