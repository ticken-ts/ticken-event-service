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
	hsm              infra.HSM
	organizerRepo    repos.OrganizerRepository
	organizationRepo repos.OrganizationRepository
}

func NewOrgManager(hsm infra.HSM, organizerRepo repos.OrganizerRepository, organizationRepo repos.OrganizationRepository) *OrgManager {
	return &OrgManager{
		hsm:              hsm,
		organizerRepo:    organizerRepo,
		organizationRepo: organizationRepo,
	}
}

func (organizationManager *OrgManager) RegisterOrganizer(organizerID string, username string, email string) (*models.Organizer, error) {
	orgWithSameID := organizationManager.organizerRepo.FindOrganizer(organizerID)
	if orgWithSameID != nil {
		return nil, fmt.Errorf("organizer %s already registerd", organizerID)
	}

	// the data comes from a JWT signed by the identity provider. This ensures that
	// the organizer is unique. The only thing that we should consider is the fact
	// that the organizer is not already registered in this server. The uniqueness
	// of the email and the username is already guaranteed by the identity provider

	organizer := models.NewOrganizer(organizerID, username, email)

	err := organizationManager.organizerRepo.AddOrganizer(organizer)
	if err != nil {
		return nil, err
	}

	return organizer, nil
}

func (organizationManager *OrgManager) RegisterOrganization(name string, organizerID string, username string) (*models.Organization, error) {
	admin := organizationManager.organizerRepo.FindOrganizerByUsername(username)
	if admin == nil {
		return nil, fmt.Errorf("organizer %s is not registerd", username)
	}

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

	cryptoMaterial, err := adm.AddOrganizationToChannel(name, admin.Username, "ticken-channel", newOrgTemplate)
	if err != nil {
		return nil, err
	}

	orgCaPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.OrgCA.GetPrivPEMEncodedBytes())
	tlsCaPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.TlsCA.GetPrivPEMEncodedBytes())
	adminOrgCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.AdminOrgCert.GetPrivPEMEncodedBytes())
	adminTlsCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.AdminTlsCert.GetPrivPEMEncodedBytes())

	newOrganization := models.NewOrganization(
		name,
		models.NewCertificate(cryptoMaterial.OrgCA.GetCertPEMEncodedBytes(), orgCaPrivStoreKey),
		models.NewCertificate(cryptoMaterial.TlsCA.GetCertPEMEncodedBytes(), tlsCaPrivStoreKey),
	)

	err = newOrganization.AddAdmin(
		admin,
		models.NewCertificate(cryptoMaterial.AdminOrgCert.GetCertPEMEncodedBytes(), adminOrgCertPrivStoreKey),
		models.NewCertificate(cryptoMaterial.AdminTlsCert.GetCertPEMEncodedBytes(), adminTlsCertPrivStoreKey),
	)
	if err != nil {
		return nil, err
	}

	err = organizationManager.organizationRepo.AddOrganization(newOrganization)
	if err != nil {
		panic(err) // TODO - Handle this situation
	}

	return newOrganization, nil
}
