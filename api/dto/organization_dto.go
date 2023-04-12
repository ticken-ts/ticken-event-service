package dto

import "github.com/google/uuid"

type CertificateDTO struct {
	Content []byte `bson:"content"`
}

type OrganizationDTO struct {
	OrganizationID uuid.UUID           `bson:"organization_id"`
	Name           string              `bson:"name"`
	MSPID          string              `bson:"msp_id"`
	Channel        string              `bson:"channel"`
	Users          []*OrganizationUser `bson:"users"`
	Nodes          []*OrganizationNode `bson:"nodes"`
}

type OrganizationUser struct {
	OrganizerID uuid.UUID `bson:"organizer_id"`
	Username    string    `bson:"username"`
	Role        string    `bson:"role"`
}

type OrganizationNode struct {
	NodeName    string          `bson:"node_name"`
	Address     string          `bson:"address"`
	NodeOrgCert *CertificateDTO `bson:"org_cert"`
	NodeTlsCert *CertificateDTO `bson:"tls_cert"`
}
