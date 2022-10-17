package services

import "ticken-event-service/models"

type IProvider interface {
	GetEventManager() EventManager
	GetUserManager() UserManager
	GetOrgManager() OrganizationManager
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
	GetEvent(eventId string, userId string) (*models.Event, error)
	GetUserEvents(userId string) ([]*models.Event, error)
	UpdateEvent(EventID string, OrganizerID string, PvtBCChannel string, Sections []models.Section) (*models.Event, error)
}

type UserManager interface {
	GetUserIdFromToken(token string) (string, error)
}

type OrganizationManager interface {
	GetUserOrganization(userId string) (*models.Organization, error)
	AddOrganization(id string, peers []string, users []string) (*models.Organization, error)
}
