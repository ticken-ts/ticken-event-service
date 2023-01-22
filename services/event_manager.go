package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/async"
	"ticken-event-service/exception"
	"ticken-event-service/log"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"time"
)

type EventManager struct {
	publisher           *async.Publisher
	eventRepo           repos.EventRepository
	organizerRepo       repos.OrganizerRepository
	organizationRepo    repos.OrganizationRepository
	organizationManager IOrganizationManager
}

func NewEventManager(
	eventRepo repos.EventRepository,
	organizerRepo repos.OrganizerRepository,
	organizationRepo repos.OrganizationRepository,
	publisher *async.Publisher,
	organizationManager IOrganizationManager,
) IEventManager {
	return &EventManager{
		publisher:           publisher,
		eventRepo:           eventRepo,
		organizerRepo:       organizerRepo,
		organizationRepo:    organizationRepo,
		organizationManager: organizationManager,
	}
}

func (eventManager *EventManager) CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time) (*models.Event, error) {
	organizer := eventManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, exception.WithMessage("organizer with id %s not found", organizerID)
	}

	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, exception.WithMessage("organization with id %s not found", organizationID)
	}

	event, err := models.NewEvent(name, date, organizer, organization)
	if err != nil {
		return nil, exception.FromError(err, "failed to create event")
	}

	if err = eventManager.eventRepo.AddEvent(event); err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	_, err = atomicPvtbcCaller.TickenEventCaller.CreateEvent(event.EventID, event.Name, event.Date)
	if err != nil {
		return nil, err
	}

	event.SetOnChain(organization.Channel)

	if err := eventManager.publisher.PublishNewEvent(event); err != nil {
		// TODO -> how to handle
		panic(err)
	}

	return event, nil
}

func (eventManager *EventManager) AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error) {
	section := models.NewSection(name, eventID, totalTickets, ticketPrice)

	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		// todo - how to handle this?
		// this case is more complicated. we should let pass
		// adding a section without the event? we should had maybe
		// some way to try to sync here with the on chain event
		return section, nil
	}

	if err := event.AssociateSection(section); err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	eventManager.eventRepo.UpdateEvent(event)

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	_, err = atomicPvtbcCaller.TickenEventCaller.AddSection(section.EventID, section.Name, section.TotalTickets, section.TicketPrice)
	if err != nil {
		return nil, err
	}

	section.OnChain = true

	return section, nil
}

func (eventManager *EventManager) GetEvent(eventID uuid.UUID, requesterID uuid.UUID, organizationID uuid.UUID) (*models.Event, error) {
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event %s not found", eventID)
	}

	requester := eventManager.organizerRepo.FindOrganizer(requesterID)
	if requester == nil {
		return nil, fmt.Errorf("requester with id %s not found", requesterID)
	}

	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("organization with id %s not found", organizationID)
	}

	if !organization.HasUser(requester.Username) {
		return nil, fmt.Errorf("user do not belong to the organization")
	}

	if !event.IsFromOrganization(organization.OrganizationID) {
		return nil, fmt.Errorf("user %s doest not belongs to the event organization", requesterID)
	}

	return event, nil
}

func (eventManager *EventManager) GetOrganizationEvents(requesterID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error) {
	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("organization with id %s not found", organizationID)
	}

	requester := eventManager.organizerRepo.FindOrganizer(requesterID)
	if requester == nil {
		return nil, fmt.Errorf("organizer with id %s not found", requesterID)
	}

	if !organization.HasUser(requester.Username) {
		return nil, fmt.Errorf("user do not belong to the organization")
	}

	return eventManager.eventRepo.FindOrganizationEvents(organizationID), nil
}
