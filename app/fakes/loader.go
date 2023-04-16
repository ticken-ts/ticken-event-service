package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"reflect"
	"strings"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/log"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

const Filename = "fakes.json"

type Loader struct {
	hsm             infra.HSM
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
}

func NewFakeLoader(repoProvider repos.IProvider, serviceProvider services.IProvider, config *config.Config, hsm infra.HSM) *Loader {
	return &Loader{
		hsm:             hsm,
		config:          config,
		repoProvider:    repoProvider,
		serviceProvider: serviceProvider,
	}
}

func (loader *Loader) Populate() error {
	if env.TickenEnv.IsProd() || !utils.FileExists(Filename) {
		return nil
	}

	seedContent := make(map[string]json.RawMessage)

	seedRawContent, err := os.ReadFile(Filename)
	if err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to read seed file: %s", err.Error()))
	}

	if err := json.Unmarshal(seedRawContent, &seedContent); err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to unmarshal seed file: %s", err.Error()))
	}

	for _, modelName := range []string{"organizer", "organization", "event"} {
		log.TickenLogger.Info().Msg(
			fmt.Sprintf("%s: %s",
				color.GreenString("seeding model"),
				color.New(color.FgBlue, color.Bold).Sprintf(modelName)),
		)

		switch modelName {

		case strings.ToLower(reflect.TypeOf(models.Event{}).Name()):
			eventsToSeed := make([]*SeedEvent, 0)

			if err := json.Unmarshal(seedContent[modelName], &eventsToSeed); err != nil {
				log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to unmarshal event values: %s", err.Error()))
				continue
			}

			seedErrors := loader.seedEvents(eventsToSeed)
			if seedErrors != nil && len(seedErrors) > 0 {
				for _, seedError := range seedErrors {
					log.TickenLogger.Error().Msg(seedError.Error())
				}
				continue
			}

		case strings.ToLower(reflect.TypeOf(models.Organizer{}).Name()):
			organizersToSeed := make([]*SeedOrganizer, 0)

			if err := json.Unmarshal(seedContent[modelName], &organizersToSeed); err != nil {
				log.TickenLogger.Error().Msg(fmt.Sprintf("failed to unmarshal organizers values: %s", err.Error()))
				continue
			}

			seedErrors := loader.seedOrganizer(organizersToSeed)
			if seedErrors != nil && len(seedErrors) > 0 {
				for _, seedError := range seedErrors {
					log.TickenLogger.Error().Msg(seedError.Error())
				}
				continue
			}

		case strings.ToLower(reflect.TypeOf(models.Organization{}).Name()):
			organizationsToSeed := make([]*SeedOrganization, 0)

			if err := json.Unmarshal(seedContent[modelName], &organizationsToSeed); err != nil {
				log.TickenLogger.Error().Msg(fmt.Sprintf("failed to unmarshal organizations values: %s", err.Error()))
				continue
			}

			seedErrors := loader.seedOrganizations(organizationsToSeed)
			if seedErrors != nil && len(seedErrors) > 0 {
				for _, seedError := range seedErrors {
					log.TickenLogger.Error().Msg(seedError.Error())
				}
				continue
			}
		}

	}
	return nil
}
