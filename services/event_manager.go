package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/api/errors"
	"ticken-event-service/async"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"time"
)

type eventManager struct {
	publisher        *async.Publisher
	eventRepo        repos.EventRepository
	organizationRepo repos.OrganizationRepository
	pvtbcConnector   *pvtbc.Caller
}

func NewEventManager(
	eventRepo repos.EventRepository,
	organizationRepo repos.OrganizationRepository,
	publisher *async.Publisher,
	pvtbcConnector *pvtbc.Caller,
) EventManager {
	return &eventManager{
		publisher:        publisher,
		eventRepo:        eventRepo,
		organizationRepo: organizationRepo,
		pvtbcConnector:   pvtbcConnector,
	}
}

func (eventManager *eventManager) CreateEvent(name string, date time.Time) (*models.Event, error) {
	event := models.NewEvent(name, date)

	// TODO -> here we need to create a new peer connection using user certificates
	err := eventManager.pvtbcConnector.TickenEventCaller.CreateAsync(event.EventID, event.Name, event.Date.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	_ = eventManager.eventRepo.AddEvent(event)
	if err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
	}

	return event, nil
}

func (eventManager *eventManager) SyncOnChainEvent(onChainEvent *chain_models.Event, channelListened string) (*models.Event, error) {
	storedEvent := eventManager.eventRepo.FindEvent(onChainEvent.EventID)

	// if the event was never seen before, in other words,
	// is not present on our database, we are going to assume
	// that it was created directly from the blockchain.
	// In this case, we are going to add the event when we
	// listen int
	if storedEvent == nil {
		newEvent := models.NewEvent(onChainEvent.Name, onChainEvent.Date)
		newEvent.EventID = onChainEvent.EventID
		err := eventManager.eventRepo.AddEvent(newEvent)
		if err != nil {
			return nil, err
		}
		storedEvent = newEvent
	}

	storedEvent.OnChain = true
	storedEvent.PvtBCChannel = channelListened
	storedEvent.OrganizationID = onChainEvent.OrganizationID

	updatedEvent := eventManager.eventRepo.UpdateEvent(storedEvent)

	err := eventManager.publisher.PublishNewEvent(updatedEvent)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (eventManager *eventManager) SyncOnChanSection(onChainSection *chain_models.Section) (*models.Event, error) {
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
		storedSection = storedEvent.AddSection(onChainSection.Name, onChainSection.TotalTickets)
	}

	storedSection.OnChain = true
	updatedEvent := eventManager.eventRepo.UpdateEvent(storedEvent)

	// TODO -> implement publishing via bus event updated

	return updatedEvent, nil
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

	if event.OrganizationID != org.OrganizationID {
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
