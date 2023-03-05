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
	DevOrgsInfo   config.Orgs
	DevEventsInfo config.Events
}

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
		return exception.WithMessage("dev user with id %s not found", populator.DevUserInfo.UserID)
	}

	organization := populator.ReposProvider.GetOrganizationRepository().FindByName(populator.DevOrgsInfo.TickenOrgName)
	if organization == nil {
		return exception.WithMessage("organization with name %s not found", populator.DevOrgsInfo.TickenOrgName)
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
		return exception.WithMessage("invalid time format: %s", populator.DevEventsInfo.EventDate)
	}

	fakeEvent := &models.Event{
		EventID:        eventID,
		Name:           populator.DevEventsInfo.EventName,
		Date:           fakeTime,
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
