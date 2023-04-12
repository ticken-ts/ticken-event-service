package mappers

import (
	"ticken-event-service/api/dto"
	"ticken-event-service/models"
)

func MapOrganizationToDTO(organization *models.Organization) *dto.OrganizationDTO {
	return &dto.OrganizationDTO{
		OrganizationID: organization.OrganizationID,
		Name:           organization.Name,
		MSPID:          organization.MSPID,
		Channel:        organization.Channel,
		Users:          MapUserListToDTO(organization.Users),
		Nodes:          MapNodeListToDTO(organization.Nodes),
	}
}

func MapOrganizationListToDTO(organizations []*models.Organization) []*dto.OrganizationDTO {
	dtos := make([]*dto.OrganizationDTO, len(organizations))
	for i, e := range organizations {
		dtos[i] = MapOrganizationToDTO(e)
	}
	return dtos
}

func MapUserToDTO(user *models.OrganizationUser) *dto.OrganizationUser {
	return &dto.OrganizationUser{
		OrganizerID: user.OrganizerID,
		Username:    user.Username,
		Role:        user.Role,
	}
}

func MapUserListToDTO(users []*models.OrganizationUser) []*dto.OrganizationUser {
	dtos := make([]*dto.OrganizationUser, len(users))
	for i, e := range users {
		dtos[i] = MapUserToDTO(e)
	}
	return dtos
}

func MapNodeToDTO(node *models.OrganizationNode) *dto.OrganizationNode {
	return &dto.OrganizationNode{
		NodeName:    node.NodeName,
		Address:     node.Address,
		NodeOrgCert: &dto.CertificateDTO{Content: node.NodeOrgCert.Content},
		NodeTlsCert: &dto.CertificateDTO{Content: node.NodeTlsCert.Content},
	}
}

func MapNodeListToDTO(nodes []*models.OrganizationNode) []*dto.OrganizationNode {
	dtos := make([]*dto.OrganizationNode, len(nodes))
	for i, e := range nodes {
		dtos[i] = MapNodeToDTO(e)
	}
	return dtos
}
