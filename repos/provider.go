package repos

import (
	"fmt"
	"ticken-event-service/config"
	"ticken-event-service/infra"
	"ticken-event-service/repos/mongoDBRepos"
)

type provider struct {
	reposFactory    IFactory
	eventRepository EventRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (IProvider, error) {
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
