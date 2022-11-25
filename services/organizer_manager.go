package services

import (
	"ticken-event-service/infra"
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
