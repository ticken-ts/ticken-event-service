package services

import (
	"fmt"
	"ticken-event-service/api/errors"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type eventManager struct {
	eventRepository        repos.EventRepository
	organizationRepository repos.OrganizationRepository
}

func NewEventManager(
	eventRepository repos.EventRepository,
	organizationRepository repos.OrganizationRepository,
) EventManager {
	newEventMan := new(eventManager)
	newEventMan.eventRepository = eventRepository
	newEventMan.organizationRepository = organizationRepository
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
	org := eventManager.organizationRepository.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf(errors.UserOrgNotFound)
	}

	event := eventManager.eventRepository.FindEvent(eventId)
	if event == nil {
		return nil, fmt.Errorf(errors.EventNotFound)
	}

	if event.OrganizerID != org.OrganizationID {
		return nil, fmt.Errorf(errors.OrgEventMismatch)
	}

	return event, nil
}

func (eventManager *eventManager) GetUserEvents(userId string) ([]*models.Event, error) {

	org := eventManager.organizationRepository.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf(errors.UserOrgNotFound)
	}

	events := eventManager.eventRepository.FindOrgEvents(org.OrganizationID)
	if events == nil {
		return nil, fmt.Errorf(errors.EventNotFound)
	}

	return events, nil
}

func (eventManager *eventManager) UpdateEvent(EventID string, OrganizerID string, PvtBCChannel string, Sections []models.Section) (*models.Event, error) {
	return eventManager.eventRepository.UpdateEvent(EventID, OrganizerID, PvtBCChannel, Sections), nil
}
