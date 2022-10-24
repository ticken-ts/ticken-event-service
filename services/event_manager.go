package services

import (
	"fmt"
	"ticken-event-service/api/errors"
	"ticken-event-service/async"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type eventManager struct {
	publisher        *async.Publisher
	eventRepo        repos.EventRepository
	organizationRepo repos.OrganizationRepository
}

func NewEventManager(eventRepo repos.EventRepository, organizationRepo repos.OrganizationRepository, publisher *async.Publisher) EventManager {
	return &eventManager{
		publisher:        publisher,
		eventRepo:        eventRepo,
		organizationRepo: organizationRepo,
	}
}

func (eventManager *eventManager) AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error) {
	event := models.NewEvent(EventID, OrganizerID, PvtBCChannel)

	err := eventManager.eventRepo.AddEvent(event)
	if err != nil {
		return nil, err
	}

	err = eventManager.publisher.PublishNewEvent(event)
	if err != nil {
		// todo -> here we should retry?
		return nil, err
	}

	return event, err
}

func (eventManager *eventManager) GetEvent(eventId string, userId string) (*models.Event, error) {
	org := eventManager.organizationRepo.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf(errors.UserOrgNotFound)
	}

	event := eventManager.eventRepo.FindEvent(eventId)
	if event == nil {
		return nil, fmt.Errorf(errors.EventNotFound)
	}

	if event.OrganizerID != org.OrganizationID {
		return nil, fmt.Errorf(errors.OrgEventMismatch)
	}

	return event, nil
}

func (eventManager *eventManager) GetUserEvents(userId string) ([]*models.Event, error) {

	org := eventManager.organizationRepo.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf(errors.UserOrgNotFound)
	}

	events := eventManager.eventRepo.FindOrgEvents(org.OrganizationID)
	if events == nil {
		return nil, fmt.Errorf(errors.EventNotFound)
	}

	return events, nil
}

func (eventManager *eventManager) UpdateEvent(EventID string, OrganizerID string, PvtBCChannel string, Sections []models.Section) (*models.Event, error) {
	return eventManager.eventRepo.UpdateEvent(EventID, OrganizerID, PvtBCChannel, Sections), nil
}
