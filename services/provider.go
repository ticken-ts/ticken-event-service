package services

import (
	"ticken-event-service/blockchain/pvtbc"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
	"ticken-event-service/utils"
)

type provider struct {
	userManager         UserManager
	eventManager        EventManager
	organizationManager OrganizationManager
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	pvtbcTickenConnector, err := pvtbc.NewConnector()
	if err != nil {
		return nil, err
	}

	provider.eventManager = NewEventManager(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		pvtbcTickenConnector,
	)

	provider.userManager = NewUserManager()

	provider.organizationManager = NewOrganizationManager()

	return provider, nil
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *provider) GetUserManager() UserManager {
	return provider.userManager
}
