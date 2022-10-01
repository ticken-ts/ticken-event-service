package services

import "ticken-event-service/models"

type Provider interface {
	GetEventManager() EventManager
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
}

type UserManager interface {
}

type OrganizationManager interface {
}