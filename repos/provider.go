package repos

import (
	"fmt"
	"ticken-event-service/config"
	"ticken-event-service/infra"
	"ticken-event-service/repos/mongoDBRepos"
)

type Provider struct {
	reposFactory        IFactory
	eventRepository     EventRepository
	organizerRepository OrganizerRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (IProvider, error) {
	provider := new(Provider)

	switch dbConfig.Driver {
	case config.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, dbConfig)
	default:
		return nil, fmt.Errorf("database driver %s not implemented", dbConfig.Driver)
	}

	return provider, nil
}

func (provider *Provider) GetEventRepository() EventRepository {
	if provider.eventRepository == nil {
		provider.eventRepository = provider.reposFactory.BuildEventRepository().(EventRepository)
	}
	return provider.eventRepository
}

func (provider *Provider) GetOrganizerRepository() OrganizerRepository {
	if provider.organizerRepository == nil {
		provider.organizerRepository = provider.reposFactory.BuildOrganizerRepository().(OrganizerRepository)
	}
	return provider.organizerRepository
}
