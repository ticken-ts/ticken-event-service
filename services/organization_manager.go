package services

import (
	"fmt"
	pvtbcadmin "github.com/ticken-ts/ticken-pvtbc-adminlib"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"os"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/utils"
)

type OrganizationManager struct {
	hsm              infra.HSM
	organizerRepo    repos.OrganizerRepository
	organizationRepo repos.OrganizationRepository
}

func NewOrganizationManager(hsm infra.HSM, organizerRepo repos.OrganizerRepository, organizationRepo repos.OrganizationRepository) *OrganizationManager {
	return &OrganizationManager{
		hsm:              hsm,
		organizerRepo:    organizerRepo,
		organizationRepo: organizationRepo,
	}
}

func (organizationManager *OrganizationManager) RegisterOrganization(name string, organizerID string, username string) (*models.Organization, error) {
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

	peerHost := "peer0." + name

	cryptoMaterial, err := adm.AddOrganizationToChannel(
		models.GenerateMspID(name),
		admin.Username,
		peerHost,
		"ticken-channel",
		newOrgTemplate,
	)
	if err != nil {
		return nil, err
	}

	orgCaPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.OrgCA.GetPrivPEMEncodedBytes())
	tlsCaPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.TlsCA.GetPrivPEMEncodedBytes())

	peerOrgCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.PeerOrgCert.GetPrivPEMEncodedBytes())
	peerTlsCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.PeerTlsCert.GetPrivPEMEncodedBytes())

	adminOrgCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.AdminOrgCert.GetPrivPEMEncodedBytes())
	adminTlsCertPrivStoreKey, _ := organizationManager.hsm.Store(cryptoMaterial.AdminTlsCert.GetPrivPEMEncodedBytes())

	newOrganization := models.NewOrganization(
		name,
		models.NewCertificate(cryptoMaterial.OrgCA.GetCertPEMEncodedBytes(), orgCaPrivStoreKey),
		models.NewCertificate(cryptoMaterial.TlsCA.GetCertPEMEncodedBytes(), tlsCaPrivStoreKey),
	)

	err = newOrganization.AddPeer(
		peerHost,
		models.NewCertificate(cryptoMaterial.PeerTlsCert.GetCertPEMEncodedBytes(), peerOrgCertPrivStoreKey),
		models.NewCertificate(cryptoMaterial.PeerTlsCert.GetCertPEMEncodedBytes(), peerTlsCertPrivStoreKey),
	)

	err = newOrganization.AddAdminMember(
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

func (organizationManager *OrganizationManager) GetPvtbcConnection(organizerID string, organizationID string) (*pvtbc.Caller, error) {
	organizer := organizationManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, fmt.Errorf("could not find organizer with ID %s", organizerID)
	}

	organization := organizationManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("could not find organization with ID %s", organizationID)
	}

	if !organization.HasMember(organizer.Username) {
		return nil, fmt.Errorf("organizer %s doesnt belong to organization %s", organizer.Username, organization.MspID)
	}

	orgMemberInfo := organization.GetMemberByUsername(organizer.Username)

	memberPrivBytes, err := organizationManager.hsm.Retrieve(orgMemberInfo.OrgCert.PrivKeyStorageKey)
	if err != nil {
		return nil, fmt.Errorf("could not get user private key from HSM: %s", err.Error())
	}

	orgMemberCert := string(orgMemberInfo.OrgCert.Content)
	orgMemberPriv := string(memberPrivBytes)

	pc := peerconnector.NewWithRawCredentials(
		organization.MspID,
		[]byte(orgMemberCert),
		[]byte(orgMemberPriv),
	)

	// todo - remove hardcode
	const peerEndpoint = "localhost:9051"
	const gatewayPeer = "peer0.org2.example.com"
	const channel = "ticken-channel"

	err = pc.ConnectWithRawTlsCert(peerEndpoint, gatewayPeer, orgMemberInfo.TlsCert.Content)
	if err != nil {
		return nil, err
	}

	pvtbcCaller, err := pvtbc.NewCaller(pc)
	if err != nil {
		return nil, err
	}

	_ = pvtbcCaller.SetChannel(channel)

	return pvtbcCaller, nil
}

func (organizationManager *OrganizationManager) GetOrganizationCryptoZipped(organizerID string, organizationID string) ([]byte, error) {
	organizer := organizationManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, fmt.Errorf("could not find organizer with ID %s", organizerID)
	}

	organization := organizationManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("could not find organization with ID %s", organizationID)
	}

	if !organization.HasMember(organizer.Username) {
		return nil, fmt.Errorf("organizer %s doesnt belong to organization %s", organizer.Username, organization.MspID)
	}

	peerFiles := make(map[string][]byte)
	for _, peer := range organization.Peers {
		peerOrgCertPriv, _ := organizationManager.hsm.Retrieve(peer.OrgCert.PrivKeyStorageKey)
		peerTlsCertPriv, _ := organizationManager.hsm.Retrieve(peer.TlsCert.PrivKeyStorageKey)

		peerFiles[peer.Host+"/msp/cacerts/"+"ca.pem"] = organization.OrgCA.Content
		peerFiles[peer.Host+"/msp/keystore/"+"priv_key"] = peerOrgCertPriv
		peerFiles[peer.Host+"/msp/signcerts/"+peer.Host+".pem"] = peer.OrgCert.Content
		peerFiles[peer.Host+"/msp/cacerts/"+"tlsca.pem"] = organization.TlsCA.Content

		peerFiles[peer.Host+"/tls/ca.crt"] = organization.TlsCA.Content
		peerFiles[peer.Host+"/tls/server.crt"] = peer.TlsCert.Content
		peerFiles[peer.Host+"/tls/server.key"] = peerTlsCertPriv
	}

	return utils.ZipFiles(peerFiles)
}
