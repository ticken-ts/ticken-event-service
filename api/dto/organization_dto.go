package dto

import "github.com/google/uuid"

type CertificateDTO struct {
	Content []byte `bson:"content"`
}

type OrganizationDTO struct {
	OrganizationID uuid.UUID           `json:"organization_id"`
	Name           string              `json:"name"`
	MSPID          string              `json:"msp_id"`
	Channel        string              `json:"channel"`
	Users          []*OrganizationUser `json:"users"`
	Nodes          []*OrganizationNode `json:"nodes"`
}

type OrganizationUser struct {
	OrganizerID uuid.UUID `json:"organizer_id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
}

type OrganizationNode struct {
	NodeName    string          `json:"node_name"`
	Address     string          `json:"address"`
	NodeOrgCert *CertificateDTO `json:"org_cert"`
	NodeTlsCert *CertificateDTO `json:"tls_cert"`
}
