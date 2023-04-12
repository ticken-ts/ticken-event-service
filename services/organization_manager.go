package services

import (
	"fmt"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/organizationerr"
	"ticken-event-service/tickenerr/organizererr"

	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
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
	pvtbcCallerAtomicBuilder pvtbcCallerAtomicBuilder,
) *OrganizationManager {
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
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}
	organization := service.organizationRepo.FindOrganization(organizationID)
	if organization == nil {
		return nil, tickenerr.New(organizationerr.OrganizationNotFoundErrorCode)
	}

	if !organization.HasUser(organizer.OrganizerID) {
		return nil, tickenerr.NewWithMessage(
			organizationerr.EstablishPVTBCConnectionErrorCode,
			fmt.Sprintf("user %s doest not belong to organization %s", organizer.Username, organization.Name),
		)
	}

	if !organization.HasNodes() {
		return nil, tickenerr.NewWithMessage(
			organizationerr.EstablishPVTBCConnectionErrorCode,
			fmt.Sprintf("organization %s doest not have any nodes", organization.Name),
		)
	}

	// we are sure that this user belongs to the organization
	// and we will find it, because we checked in the lines before
	orgUserInfo := organization.GetUserByID(organizer.OrganizerID)

	memberPrivBytes, err := service.hsm.Retrieve(orgUserInfo.UserOrgCert.PrivKeyStorageKey)
	if err != nil {
		return nil, fmt.Errorf("could not get user private key from HSM: %s", err.Error())
	}

	orgMemberPriv := memberPrivBytes
	orgMemberCert := orgUserInfo.UserOrgCert.Content

	peerNode := organization.Nodes[0]

	pvtbcCaller, err := service.pvtbcCallerAtomicBuilder(
		organization.MSPID,
		organizer.Username,
		peerNode.Address,
		orgMemberCert,
		orgMemberPriv,
		peerNode.NodeTlsCert.Content,
	)
	if err != nil {
		return nil, tickenerr.FromError(organizationerr.EstablishPVTBCConnectionErrorCode, err)
	}

	if err := pvtbcCaller.SetChannel(organization.Channel); err != nil {
		return nil, tickenerr.FromErrorWithMessage(
			organizationerr.EstablishPVTBCConnectionErrorCode,
			err,
			fmt.Sprintf("could not stablish connection in PVTBC channel %s", organization.Channel),
		)
	}

	return pvtbcCaller, nil
}

func (service *OrganizationManager) GetOrganizationsByOrganizer(organizerID uuid.UUID) ([]*models.Organization, error) {
	// check if organizer exists
	if service.organizerRepo.FindOrganizer(organizerID) == nil {
		return nil, tickenerr.New(organizererr.OrganizerNotFoundErrorCode)
	}
	return service.organizationRepo.FindByOrganizer(organizerID), nil
}
