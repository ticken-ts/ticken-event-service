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
}

type IEventManager interface {
	// CreateEvent create and event, stored it locally and publish
	// it to the private blockchain.
	// The poster (image of the event) can be set from two sources:
	// This method DO NOT publish the
	// event contract to the public blockchain, neither to publish it
	// to other services, because the event still can be modified, such
	// as adding new sections. To finally publish the event, we need
	// to call method StartSale.
	CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time, description string, poster *file.File) (*models.Event, error)

	// AddSection adds a section into the event  with "totalTickets" in it,
	// with each ticket with a price of "ticketPrice".
	// The event must not be on sale in order to be able to
	// add more sections. Each section must have a unique name
	AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error)

	GetEvent(eventID, organizerID, organizationID uuid.UUID) (*models.Event, error)
	GetOrganizationEvents(organizerID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error)

	StartSale(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error)
	StartEvent(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error)
	FinishEvent(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error)

	GetAvailableEvents() ([]*models.Event, error)
	GetPublicEvent(eventID uuid.UUID) (*models.Event, error)
}

type IOrganizationManager interface {
	GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error)
}

type IValidatorManager interface {
	RegisterValidator(organizerID, organizationID uuid.UUID, username, password, email string) (*models.Validator, error)
}

type IAssetManager interface {
	GetAssetURL(assetID uuid.UUID) (string, error)
	UploadAsset(file *file.File, name string) (*models.Asset, error)
}
