package models

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	mongoID        primitive.ObjectID    `bson:"_id"`
	OrganizationID string                `bson:"organization_id"`
	MspID          string                `bson:"msp_id"`
	OrgCA          *Certificate          `bson:"org_ca"`
	TlsCA          *Certificate          `bson:"tls_ca"`
	Members        []*OrganizationMember `bson:"members"`
}

type OrganizationMember struct {
	IsAdmin     bool         `bson:"is_admin"`
	OrganizerID string       `bson:"organizer_id"`
	Username    string       `bson:"username"`
	OrgCert     *Certificate `bson:"org_cert"`
	TlsCert     *Certificate `bson:"tls_cert"`
}

type Certificate struct {
	Content           []byte `bson:"content"`
	PrivKeyStorageKey string `bson:"private_key"`
}

func NewOrganization(orgName string, orgCA *Certificate, tlsCA *Certificate) *Organization {
	return &Organization{
		OrganizationID: uuid.New().String(),
		MspID:          orgName,
		OrgCA:          orgCA,
		TlsCA:          tlsCA,
		Members:        make([]*OrganizationMember, 0),
	}
}

func (organization *Organization) AddMember(organizer *Organizer, orgCert *Certificate, tlsCert *Certificate) error {
	if organization.HasMember(organizer.Username) {
		return fmt.Errorf("%s already is part of the organization %s", organizer.Username, organization.MspID)
	}

	newMember := &OrganizationMember{
		IsAdmin:     false,
		OrganizerID: organizer.OrganizerID,
		Username:    organizer.Username,
		OrgCert:     orgCert,
		TlsCert:     tlsCert,
	}

	organization.Members = append(organization.Members, newMember)
	return nil
}

func (organization *Organization) AddAdmin(organizer *Organizer, orgCert *Certificate, tlsCert *Certificate) error {
	if organization.HasMember(organizer.Username) {
		return fmt.Errorf("%s already is part of the organization %s", organizer.Username, organization.MspID)
	}

	newMember := &OrganizationMember{
		IsAdmin:     false,
		OrganizerID: organizer.OrganizerID,
		Username:    organizer.Username,
		OrgCert:     orgCert,
		TlsCert:     tlsCert,
	}

	organization.Members = append(organization.Members, newMember)
	return nil
}

func NewCertificate(content []byte, privKeyStorageKey string) *Certificate {
	return &Certificate{
		Content:           content,
		PrivKeyStorageKey: privKeyStorageKey,
	}
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
