package repos

import (
	"ticken-event-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
	FindOrgEvents(orgID string) []*models.Event
	UpdateEvent(EventID string, OrganizerID string, PvtBCChannel string, Sections []models.Section) *models.Event
}

type OrganizationRepository interface {
	FindUserOrganization(userId string) *models.Organization
	AddOrganization(org *models.Organization) error
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetOrganizationRepository() OrganizationRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildOrganizationRepository() any
}
