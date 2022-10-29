package main

import (
	"ticken-event-service/app"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/log"
)

func main() {
	tickenEnv, err := env.Load()
	if err != nil {
		panic(err)
	}

	log.InitGlobalLogger()

	tickenConfig, err := config.Load(tickenEnv.ConfigFilePath, tickenEnv.ConfigFileName)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}

	infraBuilder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}

	tickenEventService := app.New(infraBuilder, tickenConfig)
	if tickenEnv.IsDev() {
		tickenEventService.Populate()
		tickenEventService.EmitFakeJWT()
	}

	tickenEventService.Start()
}
