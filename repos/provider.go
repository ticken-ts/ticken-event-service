package repos

import (
	"fmt"
	"ticken-event-service/config"
	"ticken-event-service/infra"
	"ticken-event-service/repos/mongoDBRepos"
)

type provider struct {
	reposFactory           Factory
	eventRepository        EventRepository
	organizationRepository OrganizationRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (Provider, error) {
	provider := new(provider)

	switch dbConfig.Driver {
	case config.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, dbConfig)

	default:
		return nil, fmt.Errorf("database driver %s not implemented", dbConfig.Driver)
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
