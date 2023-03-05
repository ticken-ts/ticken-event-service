package models

import (
	"github.com/google/uuid"
	"ticken-event-service/exception"
	"time"
)

type Event struct {
	// ************* PVTBC Data ************* //
	EventID       uuid.UUID  `bson:"event_id"`
	Name          string     `bson:"name"`
	Date          time.Time  `bson:"date"`
	Description   string     `bson:"description"`
	Sections      []*Section `bson:"sections"`
	OnSale        bool       `bson:"on_sale"`
	PosterAssetID *uuid.UUID `bson:"poster_id"`
	// ************************************** //

	// ********** Access & Auditory ********* //
	OrganizerID    uuid.UUID `bson:"organizer_id"`
	OrganizationID uuid.UUID `bson:"organization_id"`
	// ************************************** //

	// ************ PVTBC Metadata ********** //
	OnChain      bool   `bson:"on_chain"`
	PvtBCChannel string `bson:"pvt_bc_channel"`
	// ************************************** //

	// ************ PUBBC Metadata ********** //
	PubBCAddress string `bson:"pub_bc_address"`
	// ************************************** //
}

func NewEvent(name string, date time.Time, description string, organizer *Organizer, organization *Organization) (*Event, error) {
	if !organization.HasUser(organizer.Username) {
		return nil, exception.WithMessage(
			"organizer %s doest not belong to organization %s", organizer.Username, organization.Name)
	}

	event := &Event{
		EventID:       uuid.New(),
		Name:          name,
		Date:          date,
		Description:   description,
		OnSale:        false,
		Sections:      make([]*Section, 0),
		PosterAssetID: nil,

		// this values will be validated from
		// the values that the chaincode notify us
		OrganizerID:    organizer.OrganizerID,
		OrganizationID: organization.OrganizationID,

		// on chain will become true when the
		// transaction is committed and the
		// listener updated the event state
		OnChain:      false,
		PvtBCChannel: "",
	}

	return event, nil
}

func (event *Event) SetOnChain(channel string) {
	event.OnChain = true
	event.PvtBCChannel = channel
}

func (event *Event) GetSection(name string) *Section {
	for _, section := range event.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (event *Event) AddSection(name string, totalTickets int, ticketPrice float64) *Section {
	newSection := NewSection(name, event.EventID, totalTickets, ticketPrice)
	event.Sections = append(event.Sections, newSection)
	return newSection
}

func (event *Event) AssociateSection(section *Section) {
	// just to ensure that the section has same event ID
	section.EventID = event.EventID
	event.Sections = append(event.Sections, section)
}

func (event *Event) IsFromOrganization(organizationID uuid.UUID) bool {
	return event.OrganizationID == organizationID
}
