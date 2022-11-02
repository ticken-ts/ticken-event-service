package services

import (
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/models"
	"time"
)

type IProvider interface {
	GetEventManager() IEventManager
	GetOrganizationManager() IOrganizationManager
}

type IEventManager interface {
	CreateEvent(creator string, name string, date time.Time) (*models.Event, error)
	AddSection(creator string, eventID string, name string, totalTickets int) (*models.Section, error)

	SyncOnChainEvent(onChainEvent *chain_models.Event, channelListened string) (*models.Event, error)
	SyncOnChainSection(onChainSection *chain_models.Section) (*models.Event, error)

	GetEvent(eventID string, requester string) (*models.Event, error)
	GetOrganizationEvents(requester string) ([]*models.Event, error)
}

type IOrganizationManager interface {
	RegisterOrganizer(organizerID string, name string, email string) (*models.Organizer, error)
}
