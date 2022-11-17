package services

import (
	"fmt"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type OrganizationManager struct {
	organizerRepos repos.OrganizerRepository
}

func NewOrganizationManager(organizerRepo repos.OrganizerRepository) *OrganizationManager {
	return &OrganizationManager{organizerRepos: organizerRepo}
}

func (organizationManager *OrganizationManager) RegisterOrganizer(organizerID string, username string, email string) (*models.Organizer, error) {
	orgWithSameID := organizationManager.organizerRepos.FindOrganizer(organizerID)
	if orgWithSameID != nil {
		return nil, fmt.Errorf("organizer %s already registerd", organizerID)
	}

	// the data comes from a JWT signed by the identity provider. This ensures that
	// the organizer is unique. The only thing that we should consider is the fact
	// that the organizer is not already registered in this server. The uniqueness
	// of the email and the username is already guaranteed by the identity provider

	organizer := models.NewOrganizer(organizerID, username, email)

	err := organizationManager.organizerRepos.AddOrganizer(organizer)
	if err != nil {
		return nil, err
	}

	return organizer, nil
}

func (organizationManager *OrganizationManager) RegisterOrganization(name string, organizerID string) {

}
