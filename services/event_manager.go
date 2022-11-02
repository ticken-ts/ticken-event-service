package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-event-service/async"
	"ticken-event-service/log"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/sync"
	"time"
)

type EventManager struct {
	publisher         *async.Publisher
	eventRepo         repos.EventRepository
	userServiceClient *sync.UserServiceClient
}

func NewEventManager(eventRepo repos.EventRepository, publisher *async.Publisher, userServiceClient *sync.UserServiceClient) IEventManager {
	return &EventManager{
		publisher:         publisher,
		eventRepo:         eventRepo,
		userServiceClient: userServiceClient,
	}
}

func (eventManager *EventManager) CreateEvent(creator string, name string, date time.Time) (*models.Event, error) {
	event := models.NewEvent(name, date)

	membership := eventManager.userServiceClient.GetUserMembership(creator)
	atomicPvtbcCaller, err := buildAtomicPvtbcCaller(membership)
	if err != nil {
		return nil, err
	}

	err = atomicPvtbcCaller.TickenEventCaller.CreateAsync(event.EventID, event.Name, event.Date)
	if err != nil {
		return nil, err
	}

	err = eventManager.eventRepo.AddEvent(event)
	if err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	return event, nil
}

func (eventManager *EventManager) AddSection(creator string, eventID string, name string, totalTickets int) (*models.Section, error) {
	section := models.NewSection(name, eventID, totalTickets)

	membership := eventManager.userServiceClient.GetUserMembership(creator)
	atomicPvtbcCaller, err := buildAtomicPvtbcCaller(membership)
	if err != nil {
		return nil, err
	}

	err = atomicPvtbcCaller.TickenEventCaller.AddSectionAsync(section.EventID, section.Name, section.TotalTickets)
	if err != nil {
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

	err = event.AssociateSection(section)
	if err != nil {
		// todo -> see what to do here
		// we cant fail if we couldn't save the event
		// because the tx is already submitted
		log.TickenLogger.Error().Err(err)
	}

	return section, nil
}

func (eventManager *EventManager) SyncOnChainEvent(onChainEvent *chain_models.Event, channelListened string) (*models.Event, error) {
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

	// now that we listened the event, we flagged it.
	// from this moment, this event is valid in our model
	storedEvent.SetOnChain(channelListened, onChainEvent.OrganizationID)

	updatedEvent := eventManager.eventRepo.UpdateEvent(storedEvent)

	err := eventManager.publisher.PublishNewEvent(updatedEvent)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (eventManager *EventManager) SyncOnChainSection(onChainSection *chain_models.Section) (*models.Event, error) {
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

func (eventManager *EventManager) GetEvent(eventID string, requesterID string) (*models.Event, error) {
	event := eventManager.eventRepo.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event %s not found", eventID)
	}

	requesterInfo := eventManager.userServiceClient.GetUserInfo(requesterID)

	if !event.IsFromOrganization(requesterInfo.OrganizationID) {
		return nil, fmt.Errorf("user %s doest not belongs to the event organization", requesterID)
	}

	return event, nil
}

func (eventManager *EventManager) GetOrganizationEvents(requesterID string) ([]*models.Event, error) {
	requesterInfo := eventManager.userServiceClient.GetUserInfo(requesterID)
	events := eventManager.eventRepo.FindOrgEvents(requesterInfo.OrganizationID)
	return events, nil
}

// buildAtomicPvtbcCaller creates a pvtbc caller without
// caching the grpc connection. This is useful when performing
// a transaction during a request using the user certificates
func buildAtomicPvtbcCaller(membership *sync.UserMembership) (*pvtbc.Caller, error) {
	pc := peerconnector.NewWithRawCredentials(
		membership.MspID,
		[]byte(membership.Certificate),
		[]byte(membership.PrivateKey),
	)

	err := pc.ConnectWithRawTlsCert(membership.PeerEndpoint, membership.GatewayPeer, []byte(membership.TLSCertificate))
	if err != nil {
		return nil, err
	}

	pvtbcCaller, err := pvtbc.NewCaller(pc)
	if err != nil {
		return nil, err
	}

	_ = pvtbcCaller.SetChannel("ticken-channel")

	return pvtbcCaller, nil
}
