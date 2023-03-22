package services

import (
	"fmt"
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
)

type pvtbcCallerAtomicBuilder func(mspID, user, peerAddr string, userCert, userPriv, tlsCert []byte) (*pvtbc.Caller, error)

type OrganizationManager struct {
	hsm                      infra.HSM
	organizerRepo            repos.OrganizerRepository
	organizationRepo         repos.OrganizationRepository
	pvtbcCallerAtomicBuilder pvtbcCallerAtomicBuilder
}

func NewOrganizationManager(
	repoProvider repos.IProvider,
	hsm infra.HSM,
	pvtbcCallerAtomicBuilder pvtbcCallerAtomicBuilder) *OrganizationManager {
	return &OrganizationManager{
		hsm:                      hsm,
		organizerRepo:            repoProvider.GetOrganizerRepository(),
		organizationRepo:         repoProvider.GetOrganizationRepository(),
		pvtbcCallerAtomicBuilder: pvtbcCallerAtomicBuilder,
	}
}

func (service *OrganizationManager) GetPvtbcConnection(organizerID uuid.UUID, organizationID uuid.UUID) (*pvtbc.Caller, error) {
	organizer := service.organizerRepo.FindOrganizer(organizerID)
	if organizer == nil {
		return nil, fmt.Errorf("could not find organizer with ID %s", organizerID)
	}

	organization := service.organizationRepo.FindOrganization(organizationID)
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

	memberPrivBytes, err := service.hsm.Retrieve(orgUserInfo.UserOrgCert.PrivKeyStorageKey)
	if err != nil {
		return nil, fmt.Errorf("could not get user private key from HSM: %s", err.Error())
	}

	orgMemberCert := orgUserInfo.UserOrgCert.Content
	orgMemberPriv := memberPrivBytes

	peerNode := organization.Nodes[0]

	pvtbcCaller, err := service.pvtbcCallerAtomicBuilder(
		organization.MSPID,
		organizer.Username,
		peerNode.Address,
		orgMemberCert,
		orgMemberPriv,
		peerNode.NodeTlsCert.Content,
	)

	if err := pvtbcCaller.SetChannel(organization.Channel); err != nil {
		return nil, err
	}

	return pvtbcCaller, nil
}
