package services

import (
	"ticken-event-service/async"
	"ticken-event-service/repos"
)

type provider struct {
	eventManager        EventManager
	organizationManager OrganizationManager
	userManager         *UserManager
}

func NewProvider(repoProvider repos.IProvider, publisher *async.Publisher) (IProvider, error) {
	provider := new(provider)

	eventRepo := repoProvider.GetEventRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.userManager = NewUserManager()
	provider.eventManager = NewEventManager(eventRepo, organizationRepo, publisher, provider.userManager)
	provider.organizationManager = NewOrganizationManager(eventRepo, organizationRepo)

	return provider, nil
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *provider) GetOrgManager() OrganizationManager {
	return provider.organizationManager
}
