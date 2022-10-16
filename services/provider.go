package services

import (
	"ticken-event-service/config"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
)

type provider struct {
	userManager         UserManager
	eventManager        EventManager
	organizationManager OrganizationManager
}

func NewProvider(db infra.Db, tickenConfig *config.Config) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		return nil, err
	}

	provider.eventManager = NewEventManager(
		repoProvider.GetEventRepository(),
		repoProvider.GetOrganizationRepository(),
	)

	provider.userManager = NewUserManager()

	provider.organizationManager = NewOrganizationManager(
		repoProvider.GetEventRepository(),
		repoProvider.GetOrganizationRepository(),
	)

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
