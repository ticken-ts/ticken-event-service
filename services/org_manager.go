package services

import (
	"fmt"
	pvtbcadmin "github.com/ticken-ts/ticken-pvtbc-adminlib"
	"os"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type OrgManager struct {
	hsm            infra.HSM
	organizerRepos repos.OrganizerRepository
}

func NewOrgManager(organizerRepo repos.OrganizerRepository, hsm infra.HSM) *OrgManager {
	return &OrgManager{organizerRepos: organizerRepo, hsm: hsm}
}

func (organizationManager *OrgManager) RegisterOrganizer(organizerID string, username string, email string) (*models.Organizer, error) {
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

func (organizationManager *OrgManager) RegisterOrganization(name string, organizerID string, username string) (*models.Organization, error) {
	adm, err := pvtbcadmin.New(
		pvtbcadmin.WithConfig("pvtbc.config.yaml"),
		pvtbcadmin.WithMainChannel("ticken-channel"),
		pvtbcadmin.WithMainChannelAdmin("Admin"),
		pvtbcadmin.WithTickenOrgName("org1"),
		pvtbcadmin.WithMainChannelOrderer("orderer.example.com"),
	)

	if err != nil {
		return nil, err
	}

	newOrgTemplate, err := os.ReadFile("pvtbc.neworg.json")
	if err != nil {
		return nil, err
	}

	_, err = adm.AddOrganizationToChannel(name, username, "ticken-channel", newOrgTemplate)
	if err != nil {
		return nil, err
	}

	return models.NewOrganization(), nil
}
