package services

import (
	"github.com/google/uuid"
	"ticken-event-service/exception"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type OrganizerManager struct {
	hsm              infra.HSM
	organizerRepo    repos.OrganizerRepository
	organizationRepo repos.OrganizationRepository
}

func NewOrganizerManager(hsm infra.HSM, organizerRepo repos.OrganizerRepository, organizationRepo repos.OrganizationRepository) *OrganizerManager {
	return &OrganizerManager{
		hsm:              hsm,
		organizerRepo:    organizerRepo,
		organizationRepo: organizationRepo,
	}
}

func (organizerManager *OrganizerManager) RegisterOrganizer(organizerID, firstname, lastname, username, email string) (*models.Organizer, error) {
	organizerUUID, err := uuid.Parse(organizerID)
	if err != nil {
		return nil, exception.FromError(err, "register organizer")
	}

	newOrganizer := models.NewOrganizer(organizerUUID, firstname, lastname, username, email)
	err = organizerManager.organizerRepo.AddOrganizer(newOrganizer)
	if err != nil {
		return nil, exception.FromError(err, "register organizer")
	}

	return newOrganizer, nil
}
