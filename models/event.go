package models

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	// EventStatusDraft is the status of an event that is not yet published
	EventStatusDraft EventStatus = "draft"

	// EventStatusOnSale is the status of an event that is published for sale
	EventStatusOnSale EventStatus = "on_sale"

	// EventStatusRunning is the status of an event that is currently happening
	EventStatusRunning EventStatus = "running"

	// EventStatusFinished is the status of an event that has finished
	EventStatusFinished EventStatus = "finished"
)

type Event struct {
	// ************* PVTBC Data ************* //
	EventID     uuid.UUID   `bson:"event_id"`
	Name        string      `bson:"name"`
	Date        time.Time   `bson:"date"`
	Description string      `bson:"description"`
	Sections    []*Section  `bson:"sections"`
	Status      EventStatus `bson:"status"`
	// ************************************** //

	// ************** MetaData ************** //
	PosterAssetID uuid.UUID `bson:"poster_id"`
	// ************************************** //

	// ********** Access & Auditory ********* //
	OrganizerID    uuid.UUID `bson:"organizer_id"`
	OrganizationID uuid.UUID `bson:"organization_id"`
	// ************************************** //

	// ************ PVTBC Metadata ********** //
	PvtBCTxID    string `bson:"pvtbc_tx_id"`
	PvtBCChannel string `bson:"pvtbc_channel"`
	// ************************************** //

	// ************ PUBBC Metadata ********** //
	PubBCTxID    string `bson:"pubbc_tx_id"`
	PubBCAddress string `bson:"pubbc_address"`
	// ************************************** //
}

func (event *Event) GetSection(name string) *Section {
	for _, section := range event.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (event *Event) AddSection(section *Section) {
	// just to ensure that the section has same event AssetID
	section.EventID = event.EventID
	event.Sections = append(event.Sections, section)
}

func (event *Event) IsFromOrganization(organizationID uuid.UUID) bool {
	return event.OrganizationID == organizationID
}

func (event *Event) HasPoster() bool {
	return event.PosterAssetID != uuid.Nil
}

func (event *Event) StartSale(pubbcAddr, pubbcTxID string) {
	event.PvtBCTxID = pubbcTxID
	event.PubBCAddress = pubbcAddr
	event.Status = EventStatusOnSale
}

func (event *Event) Start() {
	event.Status = EventStatusRunning
}
