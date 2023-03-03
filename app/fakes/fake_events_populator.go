package fakes

import (
	"github.com/google/uuid"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/exception"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"time"
)

type FakeEventsPopulator struct {
	ReposProvider repos.IProvider
	DevUserInfo   config.DevUser
}

func (populator *FakeEventsPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	eventID := uuid.MustParse("8709adbb-0504-4707-9cb2-867126c8172f")
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
		return exception.WithMessage("dev user with id %s not found", populator.DevUserInfo.UserID)
	}

	fakeSection := &models.Section{
		Name:         "Campo VIP",
		EventID:      eventID,
		OnChain:      true,
		TicketPrice:  100,
		TotalTickets: 100,
	}

	fakeEvent := &models.Event{
		EventID:        eventID,
		Name:           "test-event-01",
		Date:           time.Time{},
		Sections:       []*models.Section{fakeSection},
		OnSale:         true,
		OrganizerID:    uuid.New(),
		OrganizationID: organizer.OrganizerID,
		OnChain:        true,
		PvtBCChannel:   "ticken-event-name",
		PubBCAddress:   "0xfafafa",
	}

	return populator.ReposProvider.GetEventRepository().AddEvent(fakeEvent)
}
