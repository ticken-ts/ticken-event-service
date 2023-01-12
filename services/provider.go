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

func NewProvider(repoProvider repos.IProvider, busPublisher infra.BusPublisher, hsm infra.HSM) (IProvider, error) {
	provider := new(provider)

	publisher, err := async.NewPublisher(busPublisher)
	if err != nil {
		return nil, err
	}

	eventRepo := repoProvider.GetEventRepository()
	organizerRepo := repoProvider.GetOrganizerRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.organizationManager = NewOrganizationManager(hsm, organizerRepo, organizationRepo)
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
