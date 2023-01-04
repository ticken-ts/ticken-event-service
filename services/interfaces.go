package services

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/models"
	"time"
)

type IProvider interface {
	GetEventManager() IEventManager
	GetOrganizerManager() IOrganizerManager
	GetOrganizationManager() IOrganizationManager
}

type IEventManager interface {
	CreateEvent(organizerID string, organizationID string, name string, date time.Time) (*models.Event, error)
	AddSection(organizerID string, organizationID string, eventID string, name string, totalTickets int) (*models.Section, error)

	SyncOnChainEvent(onChainEvent *chain_models.Event, channelListened string) (*models.Event, error)
	SyncOnChainSection(onChainSection *chain_models.Section) (*models.Event, error)

	GetEvent(eventID string, requesterID string) (*models.Event, error)
	GetOrganizationEvents(requesterID string, organizationID string) ([]*models.Event, error)
}

type IOrganizationManager interface {
	RegisterOrganization(name string, organizerID string) (*models.Organization, error)
	GetOrganizationCryptoZipped(organizerID string, organizationID string) ([]byte, error)
	GetPvtbcConnection(organizerID string, organizationID string) (*pvtbc.Caller, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error)
}
