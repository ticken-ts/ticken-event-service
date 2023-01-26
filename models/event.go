package models

import (
	"github.com/google/uuid"
	"ticken-event-service/exception"
	"time"
)

type Event struct {
	// ************* PVTBC Data ************* //
	EventID  uuid.UUID  `bson:"event_id"`
	Name     string     `bson:"name"`
	Date     time.Time  `bson:"date"`
	Sections []*Section `bson:"sections"`
	OnSale   bool       `bson:"on_sale"`
	// ************************************** //

	// ********** Access & Auditory ********* //
	OrganizerID    uuid.UUID `bson:"organizer_id"`
	OrganizationID uuid.UUID `bson:"organization_id"`
	// ************************************** //

	// ************ PVTBC Metadata ********** //
	OnChain      bool   `bson:"on_chain"`
	PvtBCChannel string `bson:"pvt_bc_channel"`
	// ************************************** //
}

func NewEvent(name string, date time.Time, organizer *Organizer, organization *Organization) (*Event, error) {
	if !organization.HasUser(organizer.Username) {
		return nil, exception.WithMessage("organizer %s doest not belong to organization %s",
			organizer.Username, organization.Name)
	}

	event := &Event{
		EventID:  uuid.New(),
		Name:     name,
		Date:     date,
		Sections: make([]*Section, 0),

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

func (event *Event) AssociateSection(section *Section) error {
	if section.EventID != event.EventID {
		return exception.WithMessage("section does not belongs to event")
	}

	sectionWithSameName := event.GetSection(section.Name)
	if sectionWithSameName != nil {
		return exception.WithMessage(
			"section with name %s already exists in event %s", sectionWithSameName.Name, event.EventID)
	}

	event.Sections = append(event.Sections, section)
	return nil
}

func (event *Event) IsFromOrganization(organizationID uuid.UUID) bool {
	return event.OrganizationID == organizationID
}
