package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Section struct {
	Name         string `json:"name" bson:"name"`
	TotalTickets int    `json:"total_tickets" bson:"total_tickets"`
	OnChain      bool   `json:"on_chain" bson:"on_chain"`
}

type Event struct {
	mongoID primitive.ObjectID `bson:"_id"`

	EventID        string    `json:"event_id" bson:"event_id"`
	Name           string    `json:"name" bson:"name"`
	Date           time.Time `json:"date" bson:"date"`
	Sections       []Section `json:"sections" bson:"sections"`
	OrganizationID string    `json:"organization_id" bson:"organization_id"`

	OnChain      bool   `json:"on_chain" bson:"on_chain"`
	PvtBCChannel string `json:"pvt_bc_channel" bson:"pvt_bc_channel"`
}

func NewEvent(name string, date time.Time) *Event {
	return &Event{
		EventID:  uuid.New().String(),
		Name:     name,
		Date:     date,
		Sections: make([]Section, 0),

		// on chain will become true when the
		// transaction is committed and the
		// listener updated the event state
		OnChain: false,
	}
}

func (event *Event) GetSection(name string) *Section {
	for _, section := range event.Sections {
		if section.Name == name {
			return &section
		}
	}
	return nil
}

func (event *Event) AddSection(name string, totalTickets int) *Section {
	newSection := Section{
		Name:         name,
		TotalTickets: totalTickets,
	}

	event.Sections = append(event.Sections, newSection)
	return &newSection
}
