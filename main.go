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
		log.TickenLogger.Panic().Msg(err.Error())
	}

	infraBuilder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}

	tickenEventService := app.New(infraBuilder, tickenConfig)
	tickenEventService.Populate()

	if tickenEnv.IsDev() {
		tickenEventService.EmitFakeJWT()
	}

	tickenEventService.Start()
}
