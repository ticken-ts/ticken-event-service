package fakes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/services"
	"ticken-event-service/utils/file"
	"time"

	"github.com/google/uuid"
)

type FakeEventsPopulator struct {
	ServiceProvider services.IProvider
	ReposProvider   repos.IProvider
	DevUserInfo     config.DevUser
	DevOrgsInfo     config.Orgs
	DevEventsInfo   config.Events
}

func (populator *FakeEventsPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	events := populator.ReposProvider.GetEventRepository().FindAvailableEvents()
	for _, event := range events {
		if event.Name == populator.DevEventsInfo.EventName {
			return nil
		}
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

	var posterUri = populator.DevEventsInfo.EventPosterUri
	poster, err := downloadFile(posterUri)
	if err != nil {
		return err
	}

	fakeTime, err := time.Parse(time.RFC3339, populator.DevEventsInfo.EventDate)
	if err != nil {
		return fmt.Errorf("invalid time format: %s", populator.DevEventsInfo.EventDate)
	}

	var manager = populator.ServiceProvider.GetEventManager()

	fakeEvent, err := manager.CreateEvent(
		organizer.OrganizerID,
		organization.OrganizationID,
		populator.DevEventsInfo.EventName,
		fakeTime,
		populator.DevEventsInfo.EventDescription,
		poster,
	)
	if err != nil {
		return fmt.Errorf("failed to create fake event %s", err.Error())
	}

	var fakeSections []*models.Section
	for _, fakeSection := range populator.DevEventsInfo.EventSections {
		fakeSection := &models.Section{
			EventID:      fakeEvent.EventID,
			OnChain:      true,
			Name:         fakeSection.SectionName,
			TotalTickets: fakeSection.SectionQuantity,
			TicketPrice:  fakeSection.SectionPrice,
		}
		fakeSections = append(fakeSections, fakeSection)
	}

	for _, section := range fakeSections {
		_, err = manager.AddSection(
			organizer.OrganizerID,
			organization.OrganizationID,
			section.EventID,
			section.Name,
			section.TotalTickets,
			section.TicketPrice,
		)
		if err != nil {
			return fmt.Errorf("failed to add fake section %s", err.Error())
		}
	}

	_, err = populator.ServiceProvider.GetEventManager().SetEventOnSale(
		fakeEvent.EventID,
		organization.OrganizationID,
		organizer.OrganizerID,
	)
	if err != nil {
		return fmt.Errorf("failed to set fake event on sale %s", err.Error())
	}

	return nil
}

func downloadFile(URL string) (*file.File, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	imageContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &file.File{
		Content:  imageContent,
		MimeType: response.Header.Get("Content-Type"),
	}, nil
}
