package repos

import (
	"ticken-event-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
	UpdateEvent(event *models.Event) *models.Event
	FindOrganizationEvents(orgID string) []*models.Event
}

type OrganizerRepository interface {
	AddOrganizer(organizer *models.Organizer) error
	FindOrganizer(organizerID string) *models.Organizer
	FindOrganizerByUsername(username string) *models.Organizer
}

type OrganizationRepository interface {
	AddOrganization(organization *models.Organization) error
	FindOrganization(organizationID string) *models.Organization
	FindOrganizationByMspID(mspID string) *models.Organization
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetOrganizerRepository() OrganizerRepository
	GetOrganizationRepository() OrganizationRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildOrganizerRepository() any
	BuildOrganizationRepository() any
}
