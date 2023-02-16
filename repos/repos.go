package repos

import (
	"github.com/google/uuid"
	"ticken-event-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID uuid.UUID) *models.Event
	UpdateEvent(event *models.Event) *models.Event
	FindOrganizationEvents(organizationID uuid.UUID) []*models.Event
	FindAvailableEvents() []*models.Event
}

type OrganizerRepository interface {
	AnyWithID(organizerID uuid.UUID) bool
	AddOrganizer(organizer *models.Organizer) error
	FindOrganizer(organizerID uuid.UUID) *models.Organizer
	FindOrganizerByUsername(username string) *models.Organizer
}

type OrganizationRepository interface {
	AnyWithName(name string) bool
	AnyWithID(organizationID uuid.UUID) bool
	AddOrganization(organization *models.Organization) error
	FindOrganization(organizationID uuid.UUID) *models.Organization
	FindByMSPID(mspID string) *models.Organization
	FindByName(name string) *models.Organization
}

type AssetRepository interface {
	FindByID(assetID uuid.UUID) *models.Asset
	AddAsset(asset *models.Asset) error
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetOrganizerRepository() OrganizerRepository
	GetOrganizationRepository() OrganizationRepository
	GetAssetRepository() AssetRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildOrganizerRepository() any
	BuildOrganizationRepository() any
	BuildAssetRepository() any
}
