package services

import (
	"ticken-event-service/async"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
	"ticken-event-service/sync"
)

type provider struct {
	eventManager        IEventManager
	organizationManager IOrgManager
}

func NewProvider(repoProvider repos.IProvider, publisher *async.Publisher, userServiceClient *sync.UserServiceClient, hsm infra.HSM) (IProvider, error) {
	provider := new(provider)

	eventRepo := repoProvider.GetEventRepository()
	organizerRepo := repoProvider.GetOrganizerRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.eventManager = NewEventManager(eventRepo, publisher, userServiceClient)
	provider.organizationManager = NewOrgManager(hsm, organizerRepo, organizationRepo)

	return provider, nil
}

func (provider *provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *provider) GetOrgManager() IOrgManager {
	return provider.organizationManager
}
