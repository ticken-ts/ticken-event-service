package models

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	mongoID primitive.ObjectID `bson:"_id"`

	EventID        string     `json:"event_id" bson:"event_id"`
	Name           string     `json:"name" bson:"name"`
	Date           time.Time  `json:"date" bson:"date"`
	Sections       []*Section `json:"sections" bson:"sections"`
	OrganizationID string     `json:"organization_id" bson:"organization_id"`

	OnChain      bool   `json:"on_chain" bson:"on_chain"`
	PvtBCChannel string `json:"pvt_bc_channel" bson:"pvt_bc_channel"`
}

func NewEvent(name string, date time.Time) *Event {
	return &Event{
		EventID:  uuid.New().String(),
		Name:     name,
		Date:     date,
		Sections: make([]*Section, 0),

		// on chain will become true when the
		// transaction is committed and the
		// listener updated the event state
		OnChain: false,
	}
}

func (event *Event) SetOnChain(channel string, organizationID string) {
	event.OnChain = true
	event.PvtBCChannel = channel
	event.OrganizationID = organizationID
}

func (event *Event) GetSection(name string) *Section {
	for _, section := range event.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (event *Event) AddSection(name string, totalTickets int) *Section {
	newSection := NewSection(name, event.EventID, totalTickets)
	event.Sections = append(event.Sections, newSection)
	return newSection
}

func (event *Event) AssociateSection(section *Section) error {
	if section.EventID != event.EventID {
		return fmt.Errorf("section does not belongs to event")
	}

	sectionWithSameName := event.GetSection(section.Name)
	if sectionWithSameName != nil {
		return fmt.Errorf("section with name %s already exists in event %s", sectionWithSameName.Name, event.EventID)
	}

	event.Sections = append(event.Sections, section)
	return nil
}
