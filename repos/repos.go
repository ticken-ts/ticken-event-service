package repos

import (
	"ticken-event-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
}

type OrganizationRepository interface {
	FindUserOrganization(userId string) *models.Organization
}

type Provider interface {
	GetEventRepository() EventRepository
	GetOrganizationRepository() OrganizationRepository
}

type Factory interface {
	BuildEventRepository() any
	BuildOrganizationRepository() any
}
