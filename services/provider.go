package services

import (
	"ticken-event-service/async"
	"ticken-event-service/repos"
	"ticken-event-service/sync"
)

type provider struct {
	eventManager EventManager
}

func NewProvider(repoProvider repos.IProvider, publisher *async.Publisher, userServiceClient *sync.UserServiceClient) (IProvider, error) {
	provider := new(provider)

	eventRepo := repoProvider.GetEventRepository()
	provider.eventManager = NewEventManager(eventRepo, publisher, userServiceClient)

	return provider, nil
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}
