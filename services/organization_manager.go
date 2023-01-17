package services

import (
	"fmt"
	"github.com/google/uuid"
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

func (organizationManager *OrganizationManager) GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error) {
	organizer := organizationManager.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, fmt.Errorf("could not find organizer with ID %s", organizerID)
	}

	organization := organizationManager.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, fmt.Errorf("could not find organization with ID %s", organizationID)
	}

	if !organization.HasUser(organizer.Username) {
		return nil, fmt.Errorf("organizer %s doesnt belong to organization %s", organizer.Username, organization.MSPID)
	}

	if !organization.HasNodes() {
		return nil, fmt.Errorf("organization doenst has any active nodes")
	}

	orgUserInfo := organization.GetUserByName(organizer.Username)

	memberPrivBytes, err := organizationManager.hsm.Retrieve(orgUserInfo.UserOrgCert.PrivKeyStorageKey)
	if err != nil {
		return nil, fmt.Errorf("could not get user private key from HSM: %s", err.Error())
	}

	orgMemberCert := string(orgUserInfo.UserOrgCert.Content)
	orgMemberPriv := string(memberPrivBytes)

	pc := peerconnector.NewWithRawCredentials(
		organization.MSPID,
		[]byte(orgMemberCert),
		[]byte(orgMemberPriv),
	)

	peerNode := organization.Nodes[0]

	//fmt.Printf("using member cert: \n %s \n", orgMemberCert)
	//fmt.Printf("using member priv: \n %s \n", orgMemberPriv)
	//fmt.Printf("using  cert: \n %s \n", string(peerNode.NodeTlsCert.Content))
	//fmt.Printf("connectin to %s peer", organization.Name+"-"+peerNode.NodeName+".localho.st")

	err = pc.ConnectWithRawTlsCert(peerNode.Address, organization.Name+"-"+peerNode.NodeName+".localho.st", peerNode.NodeTlsCert.Content)
	if err != nil {
		return nil, err
	}

	pvtbcCaller, err := pvtbc.NewCaller(pc)
	if err != nil {
		return nil, err
	}

	_ = pvtbcCaller.SetChannel(organization.Channel)

	return pvtbcCaller, nil
}
