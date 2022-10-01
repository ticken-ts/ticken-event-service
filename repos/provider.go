package repos

import (
	"fmt"
	"ticken-event-service/infra"
	"ticken-event-service/repos/mongoDBRepos"
	"ticken-event-service/utils"
)

type provider struct {
	reposFactory           Factory
	eventRepository        EventRepository
	organizationRepository OrganizationRepository
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	switch tickenConfig.Config.Database.Driver {
	case utils.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, tickenConfig)

	default:
		return nil, fmt.Errorf("database driver %s not implemented", tickenConfig.Config.Database.Driver)
	}

	return provider, nil
}

func (provider *provider) GetEventRepository() EventRepository {
	if provider.eventRepository == nil {
		provider.eventRepository = provider.reposFactory.BuildEventRepository().(EventRepository)
	}
	return provider.eventRepository
}

func (provider *provider) GetOrganizationRepository() OrganizationRepository {
	if provider.organizationRepository == nil {
		provider.organizationRepository = provider.reposFactory.BuildOrganizationRepository().(OrganizationRepository)
	}
	return provider.organizationRepository
}
