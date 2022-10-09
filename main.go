package main

import (
	"ticken-event-service/app"
	"ticken-event-service/infra"
	"ticken-event-service/utils"
)

func main() {
	tickenConfig, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	db := builder.BuildDb()
	router := builder.BuildRouter()
	pvtbcListener := builder.BuildPvtbcListener()

	tickenEventServer := app.New(router, db, pvtbcListener, tickenConfig)
	if tickenConfig.IsDev() {
		tickenEventServer.Populate()
	}

	tickenEventServer.Start()
}
