package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-event-service/async"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
)

type provider struct {
	eventManager        IEventManager
	organizerManager    IOrganizerManager
	organizationManager IOrganizationManager
	assetManager        IAssetManager
}

func NewProvider(
	repoProvider repos.IProvider, busPublisher infra.BusPublisher, hsm infra.HSM, builder infra.IBuilder, pubbcAdmin pubbc.Admin, fileUploader infra.FileUploader,
) (IProvider, error) {
	provider := new(provider)

	publisher, err := async.NewPublisher(busPublisher)
	if err != nil {
		return nil, err
	}

	eventRepo := repoProvider.GetEventRepository()
	organizerRepo := repoProvider.GetOrganizerRepository()
	organizationRepo := repoProvider.GetOrganizationRepository()

	provider.organizationManager = NewOrganizationManager(hsm, organizerRepo, organizationRepo, builder.BuildAtomicPvtbcCaller)
	provider.eventManager = NewEventManager(eventRepo, organizerRepo, organizationRepo, publisher, provider.organizationManager, pubbcAdmin, fileUploader)
	provider.assetManager = NewAssetManager(repoProvider.GetAssetRepository())

	return provider, nil
}

func (provider *provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *provider) GetOrganizationManager() IOrganizationManager {
	return provider.organizationManager
}
