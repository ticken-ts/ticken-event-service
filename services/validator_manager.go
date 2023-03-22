package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
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
	validatorServiceClient *sync.ValidatorServiceHTTPClient) IValidatorManager {
	return &ValidatorManager{
		keycloakClient:         keycloakClient,
		authIssuer:             authIssuer,
		validatorServiceClient: validatorServiceClient,
		organizersRepo:         repoProvider.GetOrganizerRepository(),
		organizationsRepo:      repoProvider.GetOrganizationRepository(),
		validatorRepository:    repoProvider.GetValidatorRepository(),
	}
}

func (manager *ValidatorManager) RegisterValidator(organizerID, organizationID uuid.UUID, username, password, email string) (*models.Validator, error) {
	organizer := manager.organizersRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, fmt.Errorf("organizer with id %s not found", organizerID)
	}
	organization := manager.organizationsRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("organization with id %s not found", organizationID)
	}

	//if !organization.HasUser(organizer.Username) {
	//	return fmt.Errorf("user %s doest not belong to organization %s", organizer.Username, organization.Name)
	//}

	keycloakUser, err := manager.keycloakClient.RegisterUser(username, password, email)
	if err != nil {
		return nil, err
	}

	newUserJWT, err := manager.authIssuer.IssueForUser(auth.Validator, email, password)
	if err != nil {
		return nil, err
	}

	if err := manager.validatorServiceClient.RegisterValidator(organizationID, newUserJWT.Token); err != nil {
		return nil, err
	}

	validator := &models.Validator{
		ValidatorID:    keycloakUser.ID,
		Email:          keycloakUser.Email,
		CreatedBy:      organizerID,
		OrganizationID: organizationID,
	}

	if err := manager.validatorRepository.AddValidator(validator); err != nil {
		return nil, err
	}

	return validator, err
}
