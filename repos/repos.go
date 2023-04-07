package repos

import (
	"github.com/google/uuid"
	"ticken-event-service/models"
)

type BaseRepository interface {
	Count() int64
	AddOne(element any) error
	AnyWithID(id uuid.UUID) bool
}

type EventRepository interface {
	BaseRepository
	FindEvent(eventID uuid.UUID) *models.Event
	UpdateEventStatus(event *models.Event) error
	UpdatePUBBCData(event *models.Event) error
	AddSectionToEvent(eventID uuid.UUID, section *models.Section) error
	FindOrganizationEvents(organizationID uuid.UUID) []*models.Event
	FindAvailableEvents() []*models.Event
}

type OrganizerRepository interface {
	BaseRepository
	FindAll() []*models.Organizer
	FindOrganizer(organizerID uuid.UUID) *models.Organizer
	FindOrganizerByUsername(username string) *models.Organizer
}

type OrganizationRepository interface {
	BaseRepository
	AnyWithName(name string) bool
	FindOrganization(organizationID uuid.UUID) *models.Organization
	FindByMSPID(mspID string) *models.Organization
	FindByName(name string) *models.Organization
}

type AssetRepository interface {
	BaseRepository
	FindByID(assetID uuid.UUID) *models.Asset
}

type ValidatorRepository interface {
	BaseRepository
	FindValidator(validatorID uuid.UUID) *models.Validator
}

type IProvider interface {
	GetAssetRepository() AssetRepository
	GetEventRepository() EventRepository
	GetOrganizerRepository() OrganizerRepository
	GetValidatorRepository() ValidatorRepository
	GetOrganizationRepository() OrganizationRepository
}

type IFactory interface {
	BuildAssetRepository() any
	BuildEventRepository() any
	BuildOrganizerRepository() any
	BuildValidatorRepository() any
	BuildOrganizationRepository() any
}
