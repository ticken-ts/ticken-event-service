package services

import (
	"fmt"
	"github.com/google/uuid"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
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

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	if err = atomicPvtbcCaller.TickenEventCaller.CreateAsync(event.EventID, event.Name, event.Date); err != nil {
		return nil, err
	}

	if err = eventManager.eventRepo.AddEvent(event); err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	return event, nil
}

func (eventManager *EventManager) AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error) {
	section := models.NewSection(name, eventID, totalTickets, ticketPrice)

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	if err = atomicPvtbcCaller.TickenEventCaller.AddSectionAsync(section.EventID, section.Name, section.TotalTickets, section.TicketPrice); err != nil {
		return nil, err
	}

	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		// todo - how to handle this?
		// this case is more complicated. we should let pass
		// adding a section without the event? we should had maybe
		// some way to try to sync here with the on chain event
		return section, nil
	}

	if err = event.AssociateSection(section); err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	return section, nil
}

func (eventManager *EventManager) SyncOnChainEvent(onChainEvent *chainmodels.Event, channelListened string) (*models.Event, error) {
	storedEvent := eventManager.eventRepo.FindEvent(onChainEvent.EventID)

	organization := eventManager.organizationRepo.FindByMSPID(onChainEvent.MSPID)
	if organization == nil {
		return nil, exception.WithMessage("organization with MSP ID %s is not loaded", onChainEvent.MSPID)
	}

	// if the event was never seen before, in other words,
	// is not present on our database, we are going to assume
	// that it was created directly from the blockchain.
	// In this case, we are going to add the event when we listen it
	if storedEvent == nil {
		newEvent, _ := models.NewEvent(onChainEvent.Name, onChainEvent.Date, nil, nil)
		newEvent.EventID = onChainEvent.EventID
		err := eventManager.eventRepo.AddEvent(newEvent)
		if err != nil {
			return nil, err
		}
		storedEvent = newEvent
	}

	// now that we listened the event, we flagged it.
	// from this moment, this event is valid in our model
	storedEvent.SetOnChain(channelListened)

	updatedEvent := eventManager.eventRepo.UpdateEvent(storedEvent)

	err := eventManager.publisher.PublishNewEvent(updatedEvent)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (eventManager *EventManager) SyncOnChainSection(onChainSection *chainmodels.Section) (*models.Event, error) {
	storedEvent := eventManager.eventRepo.FindEvent(onChainSection.EventID)
	if storedEvent == nil {
		return nil, fmt.Errorf("event %s not founf", onChainSection.EventID)
	}

	storedSection := storedEvent.GetSection(onChainSection.Name)

	// if the section was never seen before, in other words,
	// is not present on our database, we are going to assume
	// that it was created directly from the blockchain.
	// In this case, we are going to add the section to the event
	// when we listen it from the blockchain
	if storedSection == nil {
		storedSection = storedEvent.AddSection(
			onChainSection.Name,
			onChainSection.TotalTickets,
			onChainSection.TicketPrice,
		)
	}

	storedSection.OnChain = true
	updatedEvent := eventManager.eventRepo.UpdateEvent(storedEvent)

	// TODO -> implement publishing via bus event updated

	return updatedEvent, nil
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
