package services

import (
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/organizationerr"
)

type OrganizerManager struct {
	keycloakClient      *sync.KeycloakHTTPClient
	organizersRepo      repos.OrganizerRepository
	validatorRepository repos.ValidatorRepository
	authIssuer          *auth.Issuer
}

func NewOrganizerManager(
	repoProvider repos.IProvider,
	authIssuer *auth.Issuer,
	keycloakClient *sync.KeycloakHTTPClient,
) IOrganizerManager {
	return &OrganizerManager{
		keycloakClient:      keycloakClient,
		authIssuer:          authIssuer,
		organizersRepo:      repoProvider.GetOrganizerRepository(),
		validatorRepository: repoProvider.GetValidatorRepository(),
	}
}

func (manager *OrganizerManager) RegisterOrganizer(username, password, email, firstname, lastname string) (*models.Organizer, error) {
	keycloakUser, err := manager.keycloakClient.RegisterUser(username, password, email)
	if err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	organizer := &models.Organizer{
		OrganizerID: keycloakUser.ID,
		Firstname:   firstname,
		Lastname:    lastname,
		Username:    username,
		Email:       email,
	}

	if err := manager.organizersRepo.AddOne(organizer); err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	return organizer, nil
}
