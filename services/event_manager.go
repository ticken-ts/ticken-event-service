package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-event-service/async"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/eventerr"
	organizationerr "ticken-event-service/tickenerr/organizationrerr"
	"ticken-event-service/tickenerr/organizererr"
	"ticken-event-service/utils/file"
	"time"
)

type EventManager struct {
	publisher           *async.Publisher
	eventRepo           repos.EventRepository
	organizerRepo       repos.OrganizerRepository
	organizationRepo    repos.OrganizationRepository
	organizationManager IOrganizationManager
	assetManager        IAssetManager
	pubbcAdmin          pubbc.Admin
}

func NewEventManager(
	repoProvider repos.IProvider,
	publisher *async.Publisher,
	organizationManager IOrganizationManager,
	assetManager IAssetManager,
	pubbcAdmin pubbc.Admin,
) IEventManager {
	return &EventManager{
		publisher:           publisher,
		eventRepo:           repoProvider.GetEventRepository(),
		organizerRepo:       repoProvider.GetOrganizerRepository(),
		organizationRepo:    repoProvider.GetOrganizationRepository(),
		organizationManager: organizationManager,
		assetManager:        assetManager,
		pubbcAdmin:          pubbcAdmin,
	}
}

func (eventManager *EventManager) CreateEvent(
	organizerID uuid.UUID,
	organizationID uuid.UUID,
	name string,
	date time.Time,
	description string,
	poster *file.File,
) (*models.Event, error) {

	organizer := eventManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}
	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, tickenerr.New(organizationerr.OrganizationNotFoundErrorCode)
	}

	event, err := models.NewEvent(
		name,
		date,
		description,
		organizer, // auditory
		organization,
	)
	if err != nil {
		return nil, tickenerr.FromError(eventerr.EventNotFoundErrorCode, err)
	}

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	if poster != nil {
		asset, err := eventManager.assetManager.UploadAsset(
			poster,
			fmt.Sprintf("%s-poster.%s", event.Name, poster.GetExtension()),
		)
		if err != nil {
			return nil, err
		}
		event.PosterAssetID = asset.ID
	}

	// todo -> add txID to the event?
	_, _, err = atomicPvtbcCaller.CreateEvent(event.EventID, event.Name, event.Date)
	if err != nil {
		return nil, tickenerr.FromError(eventerr.FailedToStoreEventInPVTBCErrorCode, err)
	}

	event.SetOnChain(organization.Channel)

	if err := eventManager.eventRepo.AddEvent(event); err != nil {
		return nil, tickenerr.FromError(eventerr.FailedToStoreEventInPVTBCErrorCode, err)
	}

	return event, nil
}

func (eventManager *EventManager) AddSection(
	organizerID uuid.UUID,
	organizationID uuid.UUID,
	eventID uuid.UUID,
	name string,
	totalTickets int,
	ticketPrice float64,
) (*models.Section, error) {

	section := models.NewSection(name, eventID, totalTickets, ticketPrice)

	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	event.AssociateSection(section)

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(
		organizerID,
		organizationID,
	)
	if err != nil {
		return nil, err
	}

	_, txID, err := atomicPvtbcCaller.AddSection(
		section.EventID,
		section.Name,
		section.TotalTickets,
		section.TicketPrice,
	)
	if err != nil {
		return nil, tickenerr.FromError(eventerr.FailedToAddSectionInPVTBC, err)
	}

	section.SetOnChain(txID)

	eventManager.eventRepo.UpdateEvent(event)

	return section, nil
}

func (eventManager *EventManager) GetEvent(
	eventID uuid.UUID,
	organizerID uuid.UUID,
	organizationID uuid.UUID,
) (*models.Event, error) {
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	organizer := eventManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}

	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, tickenerr.New(organizationerr.OrganizationNotFoundErrorCode)
	}

	if !event.IsFromOrganization(organization.OrganizationID) {
		return nil, tickenerr.NewWithMessage(
			eventerr.EventReadPermissionErrorCode,
			fmt.Sprintf("event doest not belongs to organization"),
		)
	}

	if !organization.HasUser(organizer.OrganizerID) {
		return nil, tickenerr.NewWithMessage(
			eventerr.EventReadPermissionErrorCode,
			fmt.Sprintf("organizer doest not belongs to organization"),
		)
	}

	return event, nil
}

func (eventManager *EventManager) GetOrganizationEvents(
	organizerID uuid.UUID,
	organizationID uuid.UUID,
) ([]*models.Event, error) {
	organizer := eventManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}

	organization := eventManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, tickenerr.New(organizationerr.OrganizationNotFoundErrorCode)
	}

	if !organization.HasUser(organizer.OrganizerID) {
		return nil, tickenerr.NewWithMessage(
			eventerr.EventReadPermissionErrorCode,
			fmt.Sprintf("organizer doest not belongs to organization"),
		)
	}

	return eventManager.eventRepo.FindOrganizationEvents(organizationID), nil
}

func (eventManager *EventManager) SetEventOnSale(
	eventID uuid.UUID,
	organizationID uuid.UUID,
	organizerID uuid.UUID,
) (*models.Event, error) {
	// this method perform all permissions checks
	event, err := eventManager.GetEvent(eventID, organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(
		organizerID,
		organizationID,
	)
	if err != nil {
		return nil, err
	}

	if _, err := atomicPvtbcCaller.SetEventOnSale(eventID); err != nil {
		return nil, tickenerr.FromError(eventerr.SetTicketOnSaleInPVTBCErrorCode, err)
	}

	addr, err := eventManager.pubbcAdmin.DeployEventContract()
	if err != nil {
		panic(err) // todo -> handle this
	}

	event.OnSale = true
	event.PubBCAddress = addr

	updatedEvent := eventManager.eventRepo.UpdateEvent(event)

	// once the event is published in the public blockchain, we sent
	// it to the other services to start commercializing  the tickets
	if err := eventManager.publisher.PublishNewEvent(event); err != nil {
		panic(err) // TODO -> how to handle
	}

	return updatedEvent, nil
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

func (eventManager *EventManager) GetPublicEvent(eventID uuid.UUID) (*models.Event, error) {
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event %s not found", eventID)
	}

	if !event.OnSale {
		return nil, fmt.Errorf("event %s is not available", eventID)
	}

	return event, nil
}
