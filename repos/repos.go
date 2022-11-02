package repos

import (
	"ticken-event-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
	FindOrgEvents(orgID string) []*models.Event
	UpdateEvent(event *models.Event) *models.Event
}

type OrganizerRepository interface {
	AddOrganizer(organizer *models.Organizer) error
	FindOrganizer(organizerID string) *models.Organizer
	FindOrganizerByUsername(username string) *models.Organizer
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetOrganizerRepository() OrganizerRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildOrganizerRepository() any
}
