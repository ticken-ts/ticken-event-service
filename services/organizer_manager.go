package services

import (
	"github.com/google/uuid"
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
	organizerID := uuid.New()

	// todo -> this is to handle standalone executions
	// how we can do this more clean and beautiful?
	// I know this is not the best way to do this, but everybody
	// knows that life is difficult and led us to some difficult
	// decisions but i want to be engineer next month :) so time
	// is something priority right now
	if manager.keycloakClient != nil {
		keycloakUser, err := manager.keycloakClient.RegisterUser(username, password, email)
		if err != nil {
			return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
		}
		organizerID = keycloakUser.ID
	}

	organizer := &models.Organizer{
		OrganizerID: organizerID,
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
