package fakes

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"time"
)

type FakeEventsPopulator struct {
	ReposProvider repos.IProvider
	DevUserInfo   config.DevUser
	DevOrgsInfo   config.Orgs
	DevEventsInfo config.Events
}

/*
   "events": {
     "event_id": "8709adbb-0504-4707-9cb2-867126c8172f",
     "event_name": "ticken event",
     "event_description": "ticken event description",
     "event_date": "2023-04-22T15:04:05Z07:00",
     "event_sections": [
       {
         "section_name": "ticken section",
         "section_price": 100,
         "section_quantity": 100
       }
     ]
   },
*/

func (populator *FakeEventsPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	eventID := uuid.MustParse(populator.DevEventsInfo.EventID)
	event := populator.ReposProvider.GetEventRepository().FindEvent(eventID)
	if event != nil {
		return nil
	}

	uuidDevUser, err := uuid.Parse(populator.DevUserInfo.UserID)
	if err != nil {
		return err
	}

	organizer := populator.ReposProvider.GetOrganizerRepository().FindOrganizer(uuidDevUser)
	if organizer == nil {
		return fmt.Errorf("dev user with id %s not found", populator.DevUserInfo.UserID)
	}

	organization := populator.ReposProvider.GetOrganizationRepository().FindByName(populator.DevOrgsInfo.TickenOrgName)
	if organization == nil {
		return fmt.Errorf("organization with name %s not found", populator.DevOrgsInfo.TickenOrgName)
	}

	fakeSections := []*models.Section{}
	for i := 0; i < len(populator.DevEventsInfo.EventSections); i++ {
		fakeSection := &models.Section{
			EventID:      eventID,
			OnChain:      true,
			Name:         populator.DevEventsInfo.EventSections[i].SectionName,
			TotalTickets: populator.DevEventsInfo.EventSections[i].SectionQuantity,
			TicketPrice:  populator.DevEventsInfo.EventSections[i].SectionPrice,
		}
		fakeSections = append(fakeSections, fakeSection)
	}

	fakeTime, err := time.Parse(time.RFC3339, populator.DevEventsInfo.EventDate)
	if err != nil {
		return fmt.Errorf("invalid time format: %s", populator.DevEventsInfo.EventDate)
	}

	fakeEvent := &models.Event{
		EventID:        eventID,
		Name:           populator.DevEventsInfo.EventName,
		Date:           fakeTime,
		Description:    populator.DevEventsInfo.EventDescription,
		Sections:       fakeSections,
		OnSale:         true,
		OrganizerID:    organizer.OrganizerID,
		OrganizationID: organization.OrganizationID,
		OnChain:        true,
		PvtBCChannel:   organization.Channel,
		PubBCAddress:   "",
	}

	return populator.ReposProvider.GetEventRepository().AddEvent(fakeEvent)
}
