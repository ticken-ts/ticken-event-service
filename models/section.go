package models

import "github.com/google/uuid"

type Section struct {
	Name         string    `bson:"name"`
	EventID      uuid.UUID `bson:"event_id"`
	TicketPrice  float64   `bson:"ticket_price"`
	TotalTickets int       `bson:"total_tickets"`

	PvtBCTxID string `bson:"pvtbc_tx_id"`
}
