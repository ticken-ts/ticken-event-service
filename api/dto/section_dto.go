package dto

type SectionDTO struct {
	EventID      string  `json:"event_id"`
	Name         string  `json:"name"`
	TotalTickets int     `json:"total_tickets"`
	Price        float64 `json:"price"`
	OnChain      bool    `json:"on_chain"`
}
