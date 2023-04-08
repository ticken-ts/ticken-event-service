package services

import (
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-event-service/models"
	"ticken-event-service/utils/file"
	"time"
)

type IProvider interface {
	GetEventManager() IEventManager
	GetOrganizationManager() IOrganizationManager
	GetAssetManager() IAssetManager
	GetValidatorManager() IValidatorManager
	GetOrganizerManager() IOrganizerManager
}

type IEventManager interface {
	CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time, description string, poster *file.File) (*models.Event, error)
	AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error)
	GetEvent(eventID, organizerID, organizationID uuid.UUID) (*models.Event, error)
	GetOrganizationEvents(organizerID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error)
	StartSale(eventID, organizerID, organizationID uuid.UUID) (*models.Event, error)
	StartEvent(eventID, organizerID, organizationID uuid.UUID) (*models.Event, error)
	FinishEvent(eventID, organizerID, organizationID uuid.UUID) (*models.Event, error)
	GetEventsOnSale(withName string, fromDate time.Time, toDate time.Time) ([]*models.Event, error)
	GetPublicEvent(eventID uuid.UUID) (*models.Event, error)
}

type IOrganizationManager interface {
	GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error)
}

type IValidatorManager interface {
	RegisterValidator(organizerID, organizationID uuid.UUID, username, password, email string) (*models.Validator, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(username, password, email, firstname, lastname string) (*models.Organizer, error)
}

type IAssetManager interface {
	GetAssetURL(assetID uuid.UUID) (string, error)
	UploadAsset(file *file.File, name string) (*models.Asset, error)
}
