package services

import (
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-event-service/models"
	"time"
)

type IProvider interface {
	GetEventManager() IEventManager
	GetOrganizationManager() IOrganizationManager
}

type IEventManager interface {
	CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time) (*models.Event, error)
	AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error)
	GetEvent(eventID uuid.UUID, requesterID uuid.UUID, organizationID uuid.UUID) (*models.Event, error)
	GetOrganizationEvents(requesterID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error)
	SetEventOnSale(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error)
}

type IOrganizationManager interface {
	GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error)
}
