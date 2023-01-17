package fakes

import (
	"github.com/google/uuid"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type FakeUsersPopulator struct {
	devUserInfo   config.DevUser
	organizerRepo repos.OrganizerRepository
}

func NewFakeUsersPopulator(organizerRepo repos.OrganizerRepository, devUserInfo config.DevUser) *FakeUsersPopulator {
	return &FakeUsersPopulator{
		devUserInfo:   devUserInfo,
		organizerRepo: organizerRepo,
	}
}

func (populator *FakeUsersPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	uuidDevUser, err := uuid.Parse(populator.devUserInfo.UserID)
	if err != nil {
		return err
	}

	if populator.organizerRepo.AnyWithID(uuidDevUser) {
		return nil
	}

	devOrganizer := models.NewOrganizer(
		uuidDevUser,
		populator.devUserInfo.Firstname,
		populator.devUserInfo.Lastname,
		populator.devUserInfo.Username,
		populator.devUserInfo.Email,
	)

	if err := populator.organizerRepo.AddOrganizer(devOrganizer); err != nil {
		return err
	}

	return nil
}
