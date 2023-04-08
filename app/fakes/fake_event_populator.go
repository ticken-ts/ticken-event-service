package fakes

import (
	"fmt"
	"ticken-event-service/utils/file"
	"time"
)

type SeedEvent struct {
	Name              string         `json:"name"`
	Date              time.Time      `json:"date"`
	Description       string         `json:"description"`
	PosterURI         string         `json:"poster_uri"`
	Sections          []*SeedSection `json:"sections"`
	OrganizerUsername string         `json:"organizer_username"`
	OrganizationName  string         `json:"organization_name"`

	// sets the event on sale after creating it
	SetOnSale bool `json:"set_on_sale"`
}

type SeedSection struct {
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
}

func (loader *Loader) seedEvents(toSeed []*SeedEvent) []error {
	var seedErrors = make([]error, 0)

	if loader.repoProvider.GetEventRepository().Count() > 0 {
		return nil
	}

	for _, event := range toSeed {
		organizer := loader.repoProvider.GetOrganizerRepository().FindOrganizerByUsername(event.OrganizerUsername)
		if organizer == nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed event %s: organizer with username %s not found", event.Name, event.OrganizerUsername),
			)
			continue
		}

		organization := loader.repoProvider.GetOrganizationRepository().FindByName(event.OrganizationName)
		if organization == nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed event %s: organization with name %s not found", event.Name, event.OrganizationName),
			)
			continue
		}

		poster, err := file.Download(event.PosterURI)
		if err != nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed event %s: could not downlaoad poster (%s)", event.Name, err.Error()),
			)
		}

		fakeEvent, err := loader.serviceProvider.GetEventManager().CreateEvent(
			organizer.OrganizerID,
			organization.OrganizationID,
			event.Name,
			event.Date,
			event.Description,
			poster,
		)
		if err != nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed event %s: %s", event.Name, err.Error()),
			)
		}

		for _, section := range event.Sections {
			_, err = loader.serviceProvider.GetEventManager().AddSection(
				organizer.OrganizerID,
				organization.OrganizationID,
				fakeEvent.EventID,
				section.Name,
				section.TotalTickets,
				section.TicketPrice,
			)

			if err != nil {
				seedErrors = append(
					seedErrors,
					fmt.Errorf("failed to seed event %s: %s", event.Name, err.Error()),
				)
			}
		}

		if event.SetOnSale {
			_, err := loader.serviceProvider.GetEventManager().StartSale(
				fakeEvent.EventID,
				organization.OrganizationID,
				organizer.OrganizerID,
			)
			if err != nil {
				seedErrors = append(
					seedErrors,
					fmt.Errorf("failed to set fake event on sale %s: %s", event.Name, err.Error()),
				)
			}
		}
	}

	return seedErrors
}
