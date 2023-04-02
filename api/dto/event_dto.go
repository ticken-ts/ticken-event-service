package dto

type EventDTO struct {
	EventID      string        `json:"event_id"`
	Name         string        `json:"name"`
	Date         string        `json:"date"`
	OnChain      bool          `json:"on_chain"`
	Sections     []*SectionDTO `json:"sections"`
	Poster       string        `json:"poster"`
	Description  string        `json:"description"`
	PubBcAddress string        `json:"pub_bc_address"`
}
