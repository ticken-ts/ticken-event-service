package services

import (
	"ticken-event-service/async"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
)

type provider struct {
	eventManager        IEventManager
	organizerManager    IOrganizerManager
	organizationManager IOrganizationManager
}

func NewProvider(repoProvider repos.IProvider, publisher *async.Publisher, hsm infra.HSM) (IProvider, error) {
	provider := new(provider)

	eventRepo := repoProvider.GetEventRepository()
	organizerRepo := repoProvider.GetOrganizerRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.eventManager = NewEventManager(eventRepo, publisher, provider.organizationManager)
	provider.organizerManager = NewOrganizerManager(hsm, organizerRepo, organizationRepo)

	return provider, nil
}

func (provider *provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *provider) GetOrganizerManager() IOrganizerManager {
	return provider.organizerManager
}

func (provider *provider) GetOrganizationManager() IOrganizationManager {
	return provider.organizationManager
}
