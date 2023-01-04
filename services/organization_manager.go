package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
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
