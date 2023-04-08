package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-event-service/async"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
)

type provider struct {
	/********* standard services *********/
	assetManager        IAssetManager
	eventManager        IEventManager
	organizerManager    IOrganizerManager
	organizationManager IOrganizationManager
	validatorManager    IValidatorManager
	/*************************************/
}

func NewProvider(
	repoProvider repos.IProvider,
	busPublisher infra.BusPublisher,
	hsm infra.HSM,
	builder infra.IBuilder,
	pubbcAdmin pubbc.Admin,
	fileUploader infra.FileUploader,
	authIssuer *auth.Issuer,
	tickenConfig *config.Config,
) (IProvider, error) {
	provider := new(provider)

	publisher, err := async.NewPublisher(busPublisher)
	if err != nil {
		return nil, err
	}

	validatorsServiceClient := sync.NewValidatorServiceHTTPClient(tickenConfig.Services.Validator, authIssuer)
	validatorsKeycloakClient := sync.NewKeycloakHTTPClient(tickenConfig.Services.Keycloak, auth.Validator, authIssuer)

	var organizersKeycloakClient *sync.KeycloakHTTPClient
	if !env.TickenEnv.IsDev() || tickenConfig.Dev.Mock.DisableAuthMock {
		organizersKeycloakClient = sync.NewKeycloakHTTPClient(tickenConfig.Services.Keycloak, auth.Organizer, authIssuer)
	}

	provider.assetManager = NewAssetManager(repoProvider.GetAssetRepository(), fileUploader)
	provider.organizerManager = NewOrganizerManager(repoProvider, authIssuer, organizersKeycloakClient)
	provider.organizationManager = NewOrganizationManager(repoProvider, hsm, builder.BuildAtomicPvtbcCaller)
	provider.validatorManager = NewValidatorManager(repoProvider, authIssuer, validatorsKeycloakClient, validatorsServiceClient)
	provider.eventManager = NewEventManager(repoProvider, publisher, provider.organizationManager, provider.assetManager, pubbcAdmin)

	return provider, nil
}

func (provider *provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *provider) GetOrganizationManager() IOrganizationManager {
	return provider.organizationManager
}

func (provider *provider) GetAssetManager() IAssetManager {
	return provider.assetManager
}

func (provider *provider) GetValidatorManager() IValidatorManager {
	return provider.validatorManager
}

func (provider *provider) GetOrganizerManager() IOrganizerManager {
	return provider.organizerManager
}
