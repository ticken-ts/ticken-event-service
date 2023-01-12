package models

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type Certificate struct {
	Content []byte `bson:"content"`

	// this property is optional, as some
	// certificate may be only the value without
	// the private key that was used to sign them
	PrivKeyStorageKey string `bson:"private_key"`
}

type Organization struct {
	OrganizationID string `bson:"organization_id"`
	Name           string `bson:"name"`
	MspID          string `bson:"msp_id"`
	Channel        string `bson:"channel"`

	OrgCACert *Certificate `bson:"org_ca"`
	TlsCACert *Certificate `bson:"tls_ca"`

	Users []*OrganizationUser `bson:"users"`
	Nodes []*OrganizationNode `bson:"nodes"`
}

type OrganizationUser struct {
	OrganizerID string       `bson:"organizer_id"`
	Username    string       `bson:"username"`
	Role        string       `bson:"role"`
	UserOrgCert *Certificate `bson:"org_cert"`
}

type OrganizationNode struct {
	NodeName    string       `bson:"node_name"`
	Address     string       `bson:"address"`
	NodeOrgCert *Certificate `bson:"org_cert"`
	NodeTlsCert *Certificate `bson:"tls_cert"`
}

func NewOrganization(name string, channel string, orgCACert *Certificate, tlsCACert *Certificate) *Organization {
	return &Organization{
		OrganizationID: uuid.New().String(),
		Name:           name,
		Channel:        channel,
		MspID:          generateMspID(name),

		OrgCACert: orgCACert,
		TlsCACert: tlsCACert,

		Users: make([]*OrganizationUser, 0),
		Nodes: make([]*OrganizationNode, 0),
	}
}

func NewCertificate(content []byte, privKeyStorageKey string) *Certificate {
	return &Certificate{
		Content:           content,
		PrivKeyStorageKey: privKeyStorageKey,
	}
}

func (organization *Organization) HasNodes() bool {
	return len(organization.Nodes) > 0
}

func (organization *Organization) AddUser(organizer *Organizer, role string, userOrgCert *Certificate) error {
	if organization.HasUser(organizer.Username) {
		return fmt.Errorf("%s already is part of the organization %s", organizer.Username, organization.Name)
	}

	newUser := &OrganizationUser{
		Role:        role,
		OrganizerID: organizer.OrganizerID,
		Username:    organizer.Username,
		UserOrgCert: userOrgCert,
	}

	organization.Users = append(organization.Users, newUser)
	return nil
}

func (organization *Organization) AddNode(nodeName string, address string, nodeOrgCert *Certificate, nodeTlsCert *Certificate) error {
	if organization.HasNode(nodeName) {
		return fmt.Errorf("%s already is part of the organization %s", nodeName, organization.Nodes)
	}

	newNode := &OrganizationNode{
		NodeName:    nodeName,
		Address:     address,
		NodeOrgCert: nodeOrgCert,
		NodeTlsCert: nodeTlsCert,
	}

	organization.Nodes = append(organization.Nodes, newNode)
	return nil
}

func (organization *Organization) HasUser(username string) bool {
	return organization.GetUserByName(username) != nil
}

func (organization *Organization) HasNode(nodeName string) bool {
	return organization.GetNodeByName(nodeName) != nil
}

func (organization *Organization) GetUserByName(username string) *OrganizationUser {
	for _, orgUser := range organization.Users {
		if orgUser.Username == username {
			return orgUser
		}
	}
	return nil
}

func (organization *Organization) GetNodeByName(nodeName string) *OrganizationNode {
	for _, orgNode := range organization.Nodes {
		if orgNode.NodeName == nodeName {
			return orgNode
		}
	}
	return nil
}

func generateMspID(orgName string) string {
	orgNameWithoutSpaces := strings.ReplaceAll(orgName, " ", "-")
	return strings.ToLower(orgNameWithoutSpaces) + "MSP"
}
