package services

import (
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/models"
	"time"
)

type IProvider interface {
	GetEventManager() EventManager
	GetOrgManager() OrganizationManager
}

type EventManager interface {
	CreateEvent(creator string, name string, date time.Time) (*models.Event, error)
	AddSection(creator string, eventID string, name string, totalTickets int) (*models.Section, error)

	SyncOnChainEvent(onChainEvent *chain_models.Event, channelListened string) (*models.Event, error)
	SyncOnChainSection(onChainSection *chain_models.Section) (*models.Event, error)

	GetEvent(eventId string, userId string) (*models.Event, error)
	GetUserEvents(userId string) ([]*models.Event, error)
}

type OrganizationManager interface {
	GetUserOrganization(userId string) (*models.Organization, error)
	AddOrganization(id string, peers []string, users []string) (*models.Organization, error)
}
