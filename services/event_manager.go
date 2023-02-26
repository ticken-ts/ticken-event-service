package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-event-service/async"
	"ticken-event-service/exception"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"time"
)

// *************+ Payloads *************** //

type CreateEventProps struct {
}

// *************+************************* //

type EventManager struct {
	publisher           *async.Publisher
	eventRepo           repos.EventRepository
	organizerRepo       repos.OrganizerRepository
	organizationRepo    repos.OrganizationRepository
	organizationManager IOrganizationManager
	pubbcAdmin          pubbc.Admin
}

func NewEventManager(
	eventRepo repos.EventRepository,
	organizerRepo repos.OrganizerRepository,
	organizationRepo repos.OrganizationRepository,
	publisher *async.Publisher,
	organizationManager IOrganizationManager,
	pubbcAdmin pubbc.Admin,
) IEventManager {
	return &EventManager{
		publisher:           publisher,
		eventRepo:           eventRepo,
		organizerRepo:       organizerRepo,
		organizationRepo:    organizationRepo,
		organizationManager: organizationManager,
		pubbcAdmin:          pubbcAdmin,
	}
}

func (eventManager *EventManager) CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time, description string, poster *models.Asset) (*models.Event, error) {
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

	if poster != nil {
		event.PosterAssetID = poster.ID
	}

	// todo -> add txID to the event?
	_, _, err = atomicPvtbcCaller.TickenEventCaller.CreateEvent(event.EventID, event.Name, event.Date)
	if err != nil {
		return nil, err
	}
	event.SetOnChain(organization.Channel)

	if err := eventManager.eventRepo.AddEvent(event); err != nil {
		return nil, exception.FromError(err, "failed to store event, please sync with the blockchain")
	}

	return event, nil
}

func (eventManager *EventManager) AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error) {
	section := models.NewSection(name, eventID, totalTickets, ticketPrice)

	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, exception.WithMessage("event %s not found, please try sync with the blockchain", eventID)
	}

	event.AssociateSection(section)

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	// todo -> add txID to the section?
	_, _, err = atomicPvtbcCaller.TickenEventCaller.AddSection(
		section.EventID, section.Name, section.TotalTickets, section.TicketPrice)
	if err != nil {
		return nil, err
	}

	section.OnChain = true
	eventManager.eventRepo.UpdateEvent(event)

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

func (eventManager *EventManager) SetEventOnSale(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error) {
	event, err := eventManager.GetEvent(eventID, organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	if _, err := atomicPvtbcCaller.SetEventOnSale(eventID); err != nil {
		return nil, err
	}

	addr, err := eventManager.pubbcAdmin.DeployEventContract()
	if err != nil {
		// TODO -> how to handle
		panic(err)
	}
	event.PubBCAddress = addr

	event.OnSale = true
	updatedEvent := eventManager.eventRepo.UpdateEvent(event)

	// once the event is published in the public blockchain, we sent
	// it to the other services to start commercializing  the tickets
	if err := eventManager.publisher.PublishNewEvent(event); err != nil {
		// TODO -> how to handle
		panic(err)
	}

	return updatedEvent, err
}

// GetAvailableEvents
// returns all the events that are on sale and have not expired
func (eventManager *EventManager) GetAvailableEvents() ([]*models.Event, error) {
	events := eventManager.eventRepo.FindAvailableEvents()
	if events == nil {
		return nil, fmt.Errorf("error querying events from database")
	}

	return events, nil
}
