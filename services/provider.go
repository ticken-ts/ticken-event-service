package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-event-service/async"
	"ticken-event-service/config"
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
	servicesConfig config.ServicesConfig,
) (IProvider, error) {
	provider := new(provider)

	publisher, err := async.NewPublisher(busPublisher)
	if err != nil {
		return nil, err
	}

	validatorsKeycloakClient := sync.NewKeycloakHTTPClient(servicesConfig.Keycloak, auth.Validator, authIssuer)
	validatorsServiceClient := sync.NewValidatorServiceHTTPClient(servicesConfig.Validator, authIssuer)

	provider.organizationManager = NewOrganizationManager(repoProvider, hsm, builder.BuildAtomicPvtbcCaller)
	provider.eventManager = NewEventManager(repoProvider, publisher, provider.organizationManager, pubbcAdmin)
	provider.assetManager = NewAssetManager(repoProvider.GetAssetRepository(), fileUploader)
	provider.validatorManager = NewValidatorManager(repoProvider, authIssuer, validatorsKeycloakClient, validatorsServiceClient)

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
