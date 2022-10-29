package models

type Section struct {
	Name         string `json:"name" bson:"name"`
	EventID      string `json:"event_id" bson:"event_id"`
	TotalTickets int    `json:"total_tickets" bson:"total_tickets"`
	OnChain      bool   `json:"on_chain" bson:"on_chain"`
}

func NewSection(name string, eventID string, totalTickets int) *Section {
	return &Section{
		Name:         name,
		EventID:      eventID,
		TotalTickets: totalTickets,
		OnChain:      false,
	}
}
