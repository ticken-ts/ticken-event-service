package services

import (
	"fmt"
	"ticken-event-service/api/errors"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type organizationManager struct {
	eventRepository        repos.EventRepository
	organizationRepository repos.OrganizationRepository
}

func NewOrganizationManager(
	eventRepository repos.EventRepository,
	organizationRepository repos.OrganizationRepository,
) OrganizationManager {
	newOrgMan := new(organizationManager)
	newOrgMan.eventRepository = eventRepository
	newOrgMan.organizationRepository = organizationRepository
	return newOrgMan
}

func (manager organizationManager) GetUserOrganization(userId string) (*models.Organization, error) {
	org := manager.organizationRepository.FindUserOrganization(userId)
	if org == nil {
		return nil, fmt.Errorf(errors.UserOrgNotFound)
	}
	return org, nil
}

func (manager organizationManager) AddOrganization(id string, peers []string, users []string) (*models.Organization, error) {
	newOrg := models.NewOrganization(id, peers, users)
	err := manager.organizationRepository.AddOrganization(newOrg)
	if err != nil {
		return nil, err
	}
	return newOrg, nil

}
