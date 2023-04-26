package services

import (
	"fmt"
	"ticken-event-service/async"
	"ticken-event-service/log"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/sync"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/commonerr"
	"ticken-event-service/tickenerr/eventerr"
	"ticken-event-service/tickenerr/organizationerr"
	"ticken-event-service/tickenerr/organizererr"
	"ticken-event-service/utils/file"
	"time"

	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
)

type EventManager struct {
	eventRepo        repos.EventRepository
	organizerRepo    repos.OrganizerRepository
	organizationRepo repos.OrganizationRepository

	publisher           *async.Publisher
	organizationManager IOrganizationManager
	assetManager        IAssetManager
	pubbcAdmin          pubbc.Admin
	pubbcCaller         pubbc.Caller

	validatorServiceClient *sync.ValidatorServiceHTTPClient
}

func NewEventManager(
	repoProvider repos.IProvider,
	publisher *async.Publisher,
	organizationManager IOrganizationManager,
	assetManager IAssetManager,
	pubbcAdmin pubbc.Admin,
	pubbcCaller pubbc.Caller,
	validatorServiceClient *sync.ValidatorServiceHTTPClient,
) IEventManager {
	return &EventManager{
		publisher:              publisher,
		eventRepo:              repoProvider.GetEventRepository(),
		organizerRepo:          repoProvider.GetOrganizerRepository(),
		organizationRepo:       repoProvider.GetOrganizationRepository(),
		organizationManager:    organizationManager,
		assetManager:           assetManager,
		pubbcAdmin:             pubbcAdmin,
		pubbcCaller:            pubbcCaller,
		validatorServiceClient: validatorServiceClient,
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

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(organizerID, organizationID)
	if err != nil {
		return nil, err
	}

	if !organization.HasUser(organizer.OrganizerID) {
		return nil, tickenerr.New(organizererr.OrganizerNotBelongsToOrganization)
	}

	var posterID = uuid.Nil
	if poster != nil {
		var assetName = fmt.Sprintf("%s-poster%s", name, poster.GetExtension())
		asset, err := eventManager.assetManager.UploadAsset(poster, assetName)
		if err != nil {
			return nil, err
		}
		posterID = asset.AssetID
	}

	eventID := uuid.New()
	_, txID, err := atomicPvtbcCaller.CreateEvent(eventID, name, date)
	if err != nil {
		return nil, tickenerr.FromError(eventerr.FailedToStoreEventInPVTBCErrorCode, err)
	}

	event := &models.Event{
		EventID:     eventID,
		Name:        name,
		Date:        date,
		Description: description,
		Status:      models.EventStatusDraft,
		Sections:    make([]*models.Section, 0),

		// metadata values or extra information
		PosterAssetID: posterID,

		OrganizerID:    organizer.OrganizerID,
		OrganizationID: organization.OrganizationID,

		// will be completed after event is co
		PvtBCChannel: organization.Channel,
		PvtBCTxID:    txID,

		// will be completed after event is on sale
		PubBCTxID:    "",
		PubBCAddress: "",
	}

	if err := eventManager.eventRepo.AddOne(event); err != nil {
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
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	atomicPvtbcCaller, err := eventManager.organizationManager.GetPvtbcConnection(
		organizerID,
		organizationID,
	)
	if err != nil {
		return nil, err
	}

	_, txID, err := atomicPvtbcCaller.AddSection(
		eventID,
		name,
		totalTickets,
		ticketPrice,
	)
	if err != nil {
		return nil, tickenerr.FromError(eventerr.FailedToAddSectionInPVTBC, err)
	}

	section := &models.Section{
		EventID:      eventID,
		TicketPrice:  ticketPrice,
		Name:         name,
		TotalTickets: totalTickets,
		PvtBCTxID:    txID,
	}

	event.AddSection(section)

	if err := eventManager.eventRepo.AddSectionToEvent(eventID, section); err != nil {
		return nil, tickenerr.FromError(commonerr.FailedToUpdateElement, err)
	}

	return section, nil
}

func (eventManager *EventManager) StartSale(
	eventID uuid.UUID,
	organizerID uuid.UUID,
	organizationID uuid.UUID,
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

	if _, err := atomicPvtbcCaller.Sell(eventID); err != nil {
		return nil, tickenerr.FromError(eventerr.StartSaleInPVTBCErrorCode, err)
	}

	addr, err := eventManager.pubbcAdmin.DeployEventContract()
	if err != nil {
		panic(err) // todo -> handle this
	}

	// TODO -> add txID
	event.StartSale(addr, "")

	_ = eventManager.eventRepo.UpdateEventStatus(event)
	_ = eventManager.eventRepo.UpdatePUBBCData(event)

	// once the event is published in the public blockchain, we sent
	// it to the other services to start commercializing  the tickets
	if err := eventManager.publisher.PublishNewEvent(event); err != nil {
		panic(err) // TODO -> how to handle
	}

	return event, nil
}

func (eventManager *EventManager) StartEvent(
	eventID uuid.UUID,
	organizerID uuid.UUID,
	organizationID uuid.UUID,
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

	if _, err := atomicPvtbcCaller.Start(eventID); err != nil {
		return nil, tickenerr.FromError(eventerr.StartEventInPVTBCErrorCode, err)
	}

	if err := eventManager.validatorServiceClient.SyncTickets(eventID); err != nil {
		panic(err) // TODO -> how to handle
	}

	event.Start()
	_ = eventManager.eventRepo.UpdateEventStatus(event)
	_ = eventManager.eventRepo.UpdatePUBBCData(event)

	if err := eventManager.publisher.PublishStatusUpdate(event); err != nil {
		panic(err) // TODO -> how to handle
	}

	return event, nil
}

func (eventManager *EventManager) FinishEvent(
	eventID uuid.UUID,
	organizerID uuid.UUID,
	organizationID uuid.UUID,
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

	if _, err := atomicPvtbcCaller.Finish(eventID); err != nil {
		return nil, tickenerr.FromError(eventerr.FinishEventInPVTBCErrorCode, err)
	}

	if _, err := eventManager.pubbcCaller.RaiseAnchors(event.PubBCAddress); err != nil {
		log.TickenLogger.Error().Msg(
			fmt.Sprintf("failed to raise event %s anchors: %s", event.EventID.String(), err.Error()))
	}

	event.Finish()
	_ = eventManager.eventRepo.UpdateEventStatus(event)
	_ = eventManager.eventRepo.UpdatePUBBCData(event)

	if err := eventManager.publisher.PublishStatusUpdate(event); err != nil {
		panic(err) // TODO -> how to handle
	}

	return event, nil
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

func (eventManager *EventManager) GetEventsOnSale(
	withName string,
	fromDate time.Time,
	toDate time.Time,
) ([]*models.Event, error) {
	events := eventManager.eventRepo.FindEvents(
		withName,
		[]models.EventStatus{models.EventStatusOnSale},
		fromDate,
		toDate,
	)
	if events == nil {
		return nil, fmt.Errorf("error querying events from database")
	}

	return events, nil
}

func (eventManager *EventManager) GetPublicEvent(
	eventID uuid.UUID,
) (*models.Event, error) {
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event %s not found", eventID)
	}

	if event.Status == models.EventStatusDraft {
		return nil, fmt.Errorf("event %s is not available", eventID)
	}

	return event, nil
}
