package fakes

import (
	"github.com/google/uuid"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/security/auth"
	"ticken-event-service/sync"
)

type FakeUsersPopulator struct {
	devConfig      config.DevConfig
	reposProvider  repos.IProvider
	keycloakClient *sync.KeycloakHTTPClient
}

func NewFakeUsersPopulator(reposProvider repos.IProvider, authIssuer *auth.Issuer, devConfig config.DevConfig, servicesConfig config.ServicesConfig) *FakeUsersPopulator {
	return &FakeUsersPopulator{
		devConfig:      devConfig,
		reposProvider:  reposProvider,
		keycloakClient: sync.NewKeycloakHTTPClient(servicesConfig.Keycloak, auth.Organizer, authIssuer),
	}
}

func (populator *FakeUsersPopulator) Populate() error {
	var uuidDevUser = uuid.MustParse(populator.devConfig.User.UserID)

	if !env.TickenEnv.IsDev() || populator.devConfig.Mock.DisableAuthMock {
		foundUserInKeycloak, _ := populator.keycloakClient.GetUserByEmail(populator.devConfig.User.Email)
		if foundUserInKeycloak != nil {
			return nil
		}

		_, err := populator.keycloakClient.RegisterUser(
			populator.devConfig.User.Username,
			populator.devConfig.User.Username,
			populator.devConfig.User.Email,
		)
		if err != nil {
			return err
		}

		foundUserInKeycloak, _ = populator.keycloakClient.GetUserByEmail(populator.devConfig.User.Email)
		uuidDevUser = foundUserInKeycloak.ID
	}

	organizerRepo := populator.reposProvider.GetOrganizerRepository()

	if organizerRepo.AnyWithID(uuidDevUser) {
		return nil
	}

	devOrganizer := models.NewOrganizer(
		uuidDevUser,
		populator.devConfig.User.Firstname,
		populator.devConfig.User.Lastname,
		populator.devConfig.User.Username,
		populator.devConfig.User.Email,
	)

	if err := organizerRepo.AddOrganizer(devOrganizer); err != nil {
		return err
	}

	return nil
}
