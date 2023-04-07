package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/organizationerr"
	"ticken-event-service/tickenerr/organizererr"
)

type ValidatorManager struct {
	keycloakClient         *sync.KeycloakHTTPClient
	validatorServiceClient *sync.ValidatorServiceHTTPClient
	organizersRepo         repos.OrganizerRepository
	organizationsRepo      repos.OrganizationRepository
	validatorRepository    repos.ValidatorRepository
	authIssuer             *auth.Issuer
}

func NewValidatorManager(
	repoProvider repos.IProvider,
	authIssuer *auth.Issuer,
	keycloakClient *sync.KeycloakHTTPClient,
	validatorServiceClient *sync.ValidatorServiceHTTPClient,
) IValidatorManager {
	return &ValidatorManager{
		keycloakClient:         keycloakClient,
		authIssuer:             authIssuer,
		validatorServiceClient: validatorServiceClient,
		organizersRepo:         repoProvider.GetOrganizerRepository(),
		organizationsRepo:      repoProvider.GetOrganizationRepository(),
		validatorRepository:    repoProvider.GetValidatorRepository(),
	}
}

func (manager *ValidatorManager) RegisterValidator(
	organizerID uuid.UUID,
	organizationID uuid.UUID,
	username string,
	password string,
	email string,
) (*models.Validator, error) {
	organizer := manager.organizersRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}
	organization := manager.organizationsRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, tickenerr.New(organizationerr.OrganizationNotFoundErrorCode)
	}

	if !organization.HasUser(organizerID) {
		return nil, tickenerr.NewWithMessage(
			organizationerr.RegisterValidatorErrorCode,
			fmt.Sprintf("user %s doest not belong to organization %s", organizer.Username, organization.Name),
		)
	}

	keycloakUser, err := manager.keycloakClient.RegisterUser(username, password, email)
	if err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	newUserJWT, err := manager.authIssuer.IssueForUser(auth.Validator, email, password)
	if err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	if err := manager.validatorServiceClient.RegisterValidator(organizationID, newUserJWT.Token); err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	validator := &models.Validator{
		ValidatorID:    keycloakUser.ID,
		Email:          keycloakUser.Email,
		CreatedBy:      organizerID,
		OrganizationID: organizationID,
	}

	if err := manager.validatorRepository.AddOne(validator); err != nil {
		return nil, tickenerr.FromError(organizationerr.RegisterValidatorErrorCode, err)
	}

	return validator, nil
}
