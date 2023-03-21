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
	GetAssetManager() IAssetManager
}

type IEventManager interface {
	CreateEvent(organizerID, organizationID uuid.UUID, name string, date time.Time, description string, poster *models.Asset) (*models.Event, error)
	AddSection(organizerID, organizationID, eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*models.Section, error)
	GetEvent(eventID uuid.UUID, requesterID uuid.UUID, organizationID uuid.UUID) (*models.Event, error)
	GetOrganizationEvents(requesterID uuid.UUID, organizationID uuid.UUID) ([]*models.Event, error)
	SetEventOnSale(eventID, organizationID, organizerID uuid.UUID) (*models.Event, error)
	GetAvailableEvents() ([]*models.Event, error)
	GetPublicEvent(eventID uuid.UUID) (*models.Event, error)
}

type IOrganizationManager interface {
	GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error)
}

type IOrganizerManager interface {
	RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error)
}

type IAssetManager interface {
	GetAsset(assetID uuid.UUID) (*models.Asset, error)
	NewAsset(name string, mimeType string, url string) (*models.Asset, error)
	UploadAsset(file *models.File, name string) (*models.Asset, error)
}
