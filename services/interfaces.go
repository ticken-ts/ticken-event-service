package services

import (
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/models"
	"time"
)

type IProvider interface {
	GetEventManager() IEventManager
	GetOrganizerManager() IOrganizerManager
	GetOrganizationManager() IOrganizationManager
}

type IEventManager interface {
	CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time) (*models.Event, error)
	AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error)

	SyncOnChainEvent(onChainEvent *chainmodels.Event, channelListened string) (*models.Event, error)
	SyncOnChainSection(onChainSection *chainmodels.Section) (*models.Event, error)

	GetEvent(eventID uuid.UUID, requesterID uuid.UUID, organizationID uuid.UUID) (*models.Event, error)
	GetOrganizationEvents(requesterID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error)
}

type IOrganizationManager interface {
	GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error)
}
