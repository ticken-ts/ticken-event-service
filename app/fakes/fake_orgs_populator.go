package fakes

import (
	"github.com/google/uuid"
	"io/ioutil"
	"path"
	"strconv"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/exception"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

const clusterStoragePath = "/tmp/ticken-pv"
const totalFakeOrganizations = 3
const tickenOrg = "ticken"

type FakeOrgsPopulator struct {
	hsm           infra.HSM
	devUserInfo   config.DevUser
	reposProvider repos.IProvider
}

func NewFakeOrgsPopulator(reposProvider repos.IProvider, devUserInfo config.DevUser, hsm infra.HSM) *FakeOrgsPopulator {
	return &FakeOrgsPopulator{
		hsm:           hsm,
		devUserInfo:   devUserInfo,
		reposProvider: reposProvider,
	}
}

func (populator *FakeOrgsPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	uuidDevUser, err := uuid.Parse(populator.devUserInfo.UserID)
	if err != nil {
		return err
	}

	organizerRepo := populator.reposProvider.GetOrganizerRepository()
	organizationRepo := populator.reposProvider.GetOrganizationRepository()

	organizer := organizerRepo.FindOrganizer(uuidDevUser)
	if organizer == nil {
		return exception.WithMessage("dev user with id %s not found", populator.devUserInfo.UserID)
	}

	// load genesis org in the database
	if !organizationRepo.AnyWithName(tickenOrg) {
		populator.createOrganization(tickenOrg, organizer)
	}

	for i := 1; i <= totalFakeOrganizations; i++ {
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
	baseMspPath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "msp")

	orgCACertBytes, err := ioutil.ReadFile(path.Join(baseMspPath, "cacerts", "ca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	tlsCACertBytes, err := ioutil.ReadFile(path.Join(baseMspPath, "tlscacerts", "tlsca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(orgCACertBytes, ""), models.NewCertificate(tlsCACertBytes, "")
}

func (populator *FakeOrgsPopulator) readOrgUserMSP(orgName string, username string) *models.Certificate {
	// todo -> replace here and in the pvtbc bootstrap the admin name
	userMspPath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "users", orgName+"-admin", "msp")

	userCertBytes, err := ioutil.ReadFile(path.Join(userMspPath, "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	userPrivBytes, err := ioutil.ReadFile(path.Join(userMspPath, "keystore", "priv.pem"))
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
	nodePath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "nodes", orgName+"-"+node)

	nodeCertBytes, err := ioutil.ReadFile(path.Join(nodePath, "msp", "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivBytes, err := ioutil.ReadFile(path.Join(nodePath, "msp", "keystore", "priv.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivStorageKey, err := populator.hsm.Store(nodePrivBytes)
	if err != nil {
		panic(err)
	}

	tlsCertBytes, err := ioutil.ReadFile(path.Join(nodePath, "tls", "signcerts", "tls-cert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(nodeCertBytes, nodePrivStorageKey), models.NewCertificate(tlsCertBytes, "")
}

func getOrgChannelName(orgName string) string {
	return tickenOrg + "-" + orgName + "-" + "channel"
}
