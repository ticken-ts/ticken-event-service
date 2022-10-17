package services

import (
	"ticken-event-service/repos"
)

type provider struct {
	userManager         UserManager
	eventManager        EventManager
	organizationManager OrganizationManager
}

func NewProvider(repoProvider repos.IProvider) (IProvider, error) {
	provider := new(provider)

	eventRepo := repoProvider.GetEventRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.userManager = NewUserManager()
	provider.eventManager = NewEventManager(eventRepo, organizationRepo)
	provider.organizationManager = NewOrganizationManager(eventRepo, organizationRepo)

	return provider, nil
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *provider) GetUserManager() UserManager {
	return provider.userManager
}

func (provider *provider) GetOrgManager() OrganizationManager {
	return provider.organizationManager
}
