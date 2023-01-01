package models

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type Organization struct {
	OrganizationID string                `bson:"organization_id"`
	Name           string                `bson:"name"`
	MspID          string                `bson:"msp_id"`
	OrgCA          *Certificate          `bson:"org_ca"`
	TlsCA          *Certificate          `bson:"tls_ca"`
	Members        []*OrganizationMember `bson:"members"`
	Peers          []*OrganizationPeer   `bson:"peers"`
}

type OrganizationMember struct {
	IsAdmin     bool         `bson:"is_admin"`
	OrganizerID string       `bson:"organizer_id"`
	Username    string       `bson:"username"`
	OrgCert     *Certificate `bson:"org_cert"`
	TlsCert     *Certificate `bson:"tls_cert"`
}

type OrganizationPeer struct {
	Host    string       `bson:"host"`
	OrgCert *Certificate `bson:"org_cert"`
	TlsCert *Certificate `bson:"tls_cert"`
}

type Certificate struct {
	Content           []byte `bson:"content"`
	PrivKeyStorageKey string `bson:"private_key"`
}

func NewOrganization(name string, orgCA *Certificate, tlsCA *Certificate) *Organization {
	return &Organization{
		OrganizationID: uuid.New().String(),
		Name:           name,
		MspID:          GenerateMspID(name),
		OrgCA:          orgCA,
		TlsCA:          tlsCA,
		Members:        make([]*OrganizationMember, 0),
	}
}

func NewCertificate(content []byte, privKeyStorageKey string) *Certificate {
	return &Certificate{
		Content:           content,
		PrivKeyStorageKey: privKeyStorageKey,
	}
}

func (organization *Organization) AddRegularMember(organizer *Organizer, orgCert *Certificate, tlsCert *Certificate) error {
	if organization.HasMember(organizer.Username) {
		return fmt.Errorf("%s already is part of the organization %s", organizer.Username, organization.MspID)
	}

	newMember := &OrganizationMember{
		IsAdmin:     false,
		OrganizerID: organizer.OrganizerID.String(),
		Username:    organizer.Username,
		OrgCert:     orgCert,
		TlsCert:     tlsCert,
	}

	organization.Members = append(organization.Members, newMember)
	return nil
}

func (organization *Organization) AddAdminMember(organizer *Organizer, orgCert *Certificate, tlsCert *Certificate) error {
	if organization.HasMember(organizer.Username) {
		return fmt.Errorf("%s already is part of the organization %s", organizer.Username, organization.MspID)
	}

	newMember := &OrganizationMember{
		IsAdmin:     false,
		OrganizerID: organizer.OrganizerID.String(),
		Username:    organizer.Username,
		OrgCert:     orgCert,
		TlsCert:     tlsCert,
	}

	organization.Members = append(organization.Members, newMember)
	return nil
}

func (organization *Organization) AddPeer(hostName string, orgCert *Certificate, tlsCert *Certificate) error {
	newPeer := &OrganizationPeer{
		Host:    hostName,
		OrgCert: orgCert,
		TlsCert: tlsCert,
	}

	organization.Peers = append(organization.Peers, newPeer)
	return nil
}

func (organization *Organization) HasMember(username string) bool {
	return organization.GetMemberByUsername(username) != nil
}

func (organization *Organization) GetMemberByUsername(username string) *OrganizationMember {
	for _, member := range organization.Members {
		if member.Username == username {
			return member
		}
	}
	return nil
}

func GenerateMspID(orgName string) string {
	orgNameWithoutSpaces := strings.ReplaceAll(orgName, " ", "-")
	return strings.ToUpper(orgNameWithoutSpaces)
}
