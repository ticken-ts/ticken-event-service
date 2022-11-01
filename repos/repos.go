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

type IProvider interface {
	GetEventRepository() EventRepository
}

type IFactory interface {
	BuildEventRepository() any
}
