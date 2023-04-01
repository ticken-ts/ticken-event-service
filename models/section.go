package models

import "github.com/google/uuid"

type Section struct {
	Name         string    `bson:"name"`
	EventID      uuid.UUID `bson:"event_id"`
	TicketPrice  float64   `bson:"ticket_price"`
	TotalTickets int       `bson:"total_tickets"`
	OnChain      bool      `bson:"on_chain"`
	TxID         string    `bson:"tx_id"`
}

func NewSection(name string, eventID uuid.UUID, totalTickets int, ticketPrice float64) *Section {
	return &Section{
		EventID:      eventID,
		TicketPrice:  ticketPrice,
		Name:         name,
		TotalTickets: totalTickets,
		OnChain:      false,
	}
}

func (section *Section) SetOnChain(txID string) {
	section.OnChain = true
	section.TxID = txID
}
