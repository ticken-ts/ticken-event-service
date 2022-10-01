package services

import (
	"fmt"
	"ticken-event-service/blockchain/pvtbc"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type eventManager struct {
	eventRepository        repos.EventRepository
	organizationRepository repos.OrganizationRepository
	pvtbcConnector         pvtbc.TickenConnector
}

func NewEventManager(
	eventRepository repos.EventRepository,
	organizationRepository repos.OrganizationRepository,
	pvtbcConnector pvtbc.TickenConnector,
) EventManager {
	newEventMan := new(eventManager)
	newEventMan.eventRepository = eventRepository
	newEventMan.organizationRepository = organizationRepository
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

	org := eventManager.organizationRepository.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf("user organization not found")
	}

	if !contains(org.Events, eventId) {
		return nil, fmt.Errorf("event not in organization")
	}

	event := eventManager.eventRepository.FindEvent(eventId)
	if event == nil {
		return nil, fmt.Errorf("event not foundx")
	}
	return event, nil
}

func contains(list []string, element string) bool {
	for _, currentElement := range list {
		if currentElement == element {
			return true
		}
	}
	return false
}
