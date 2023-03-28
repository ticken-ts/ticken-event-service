package fakes

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path"
	"strconv"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
)

type FakeOrgsPopulator struct {
	hsm                infra.HSM
	devConfig          config.DevConfig
	reposProvider      repos.IProvider
	keycloakClient     *sync.KeycloakHTTPClient
	clusterStoragePath string
}

func NewFakeOrgsPopulator(
	reposProvider repos.IProvider,
	authIssuer *auth.Issuer,
	devConfig config.DevConfig,
	hsm infra.HSM,
	clusterStoragePath string) *FakeOrgsPopulator {
	return &FakeOrgsPopulator{
		hsm:                hsm,
		devConfig:          devConfig,
		reposProvider:      reposProvider,
		clusterStoragePath: clusterStoragePath,
		keycloakClient:     sync.NewKeycloakHTTPClient("http://localhost:8080", auth.Organizer, authIssuer),
	}
}

func (populator *FakeOrgsPopulator) Populate() error {
	uuidDevUser := uuid.MustParse(populator.devConfig.User.UserID)
	if !env.TickenEnv.IsDev() || populator.devConfig.Mock.DisableAuthMock {
		foundUserInKeycloak, _ := populator.keycloakClient.GetUserByEmail(populator.devConfig.User.Email)
		if foundUserInKeycloak == nil {
			return fmt.Errorf("auth is not mocked but admin is not present in identity provider")
		}
		uuidDevUser = foundUserInKeycloak.ID
	}

	organizerRepo := populator.reposProvider.GetOrganizerRepository()
	organizationRepo := populator.reposProvider.GetOrganizationRepository()

	organizer := organizerRepo.FindOrganizer(uuidDevUser)
	if organizer == nil {
		return fmt.Errorf("dev user with id %s not found", populator.devConfig.User.UserID)
	}

	// load genesis org in the database
	if !organizationRepo.AnyWithName(populator.devConfig.Orgs.TickenOrgName) {
		populator.createOrganization(populator.devConfig.Orgs.TickenOrgName, organizer)
	}

	for i := 1; i <= populator.devConfig.Orgs.TotalFakeOrgs; i++ {
		orgName := "org" + strconv.Itoa(i)
		if organizationRepo.AnyWithName(orgName) {
			continue
		}
		populator.createOrganization(orgName, organizer)
	}

	return nil
}

func (populator *FakeOrgsPopulator) createOrganization(orgName string, admin *models.Organizer) *models.Organization {
	orgCACert, tlsCACert := populator.readOrgMSP(orgName)
	newOrganization := models.NewOrganization(orgName, getOrgChannelName(orgName), orgCACert, tlsCACert)

	userOrgCert := populator.readOrgUserMSP(orgName, admin.Username)
	_ = newOrganization.AddUser(admin, "admin", userOrgCert)

	nodeOrgCert, nodeTlsCert := populator.readOrgNodeMSP(orgName, "peer0")
	_ = newOrganization.AddNode("peer0", orgName+"-peer0"+".localho.st:443", nodeOrgCert, nodeTlsCert)

	organizationRepo := populator.reposProvider.GetOrganizationRepository()

	if err := organizationRepo.AddOrganization(newOrganization); err != nil {
		panic(err)
	}
	return newOrganization
}

func (populator *FakeOrgsPopulator) readOrgMSP(orgName string) (*models.Certificate, *models.Certificate) {
	baseMspPath := path.Join(populator.clusterStoragePath, "orgs", "peer-orgs", orgName, "msp")

	orgCACertBytes, err := os.ReadFile(path.Join(baseMspPath, "cacerts", "ca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	tlsCACertBytes, err := os.ReadFile(path.Join(baseMspPath, "tlscacerts", "tlsca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(orgCACertBytes, ""), models.NewCertificate(tlsCACertBytes, "")
}

func (populator *FakeOrgsPopulator) readOrgUserMSP(orgName string, _ string) *models.Certificate {
	// todo -> replace here and in the pvtbc bootstrap the admin name
	userMspPath := path.Join(populator.clusterStoragePath, "orgs", "peer-orgs", orgName, "users", orgName+"-admin", "msp")

	userCertBytes, err := os.ReadFile(path.Join(userMspPath, "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	userPrivBytes, err := os.ReadFile(path.Join(userMspPath, "keystore", "priv.pem"))
	if err != nil {
		panic(err)
	}

	userPrivStorageKey, err := populator.hsm.Store(userPrivBytes)
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(userCertBytes, userPrivStorageKey)
}

func (populator *FakeOrgsPopulator) readOrgNodeMSP(orgName string, node string) (*models.Certificate, *models.Certificate) {
	nodePath := path.Join(populator.clusterStoragePath, "orgs", "peer-orgs", orgName, "nodes", orgName+"-"+node)

	nodeCertBytes, err := os.ReadFile(path.Join(nodePath, "msp", "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivBytes, err := os.ReadFile(path.Join(nodePath, "msp", "keystore", "priv.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivStorageKey, err := populator.hsm.Store(nodePrivBytes)
	if err != nil {
		panic(err)
	}

	tlsCertBytes, err := os.ReadFile(path.Join(nodePath, "tls", "signcerts", "tls-cert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(nodeCertBytes, nodePrivStorageKey), models.NewCertificate(tlsCertBytes, "")
}

func getOrgChannelName(orgName string) string {
	return orgName + "-" + "channel"
}
