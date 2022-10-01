package app

import (
	"ticken-event-service/api"
	"ticken-event-service/api/controllers/eventController"
	"ticken-event-service/api/controllers/organizationController"
	"ticken-event-service/infra"
	"ticken-event-service/middlewares"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

type TickenEventApp struct {
	router          infra.Router
	serviceProvider services.Provider
}

func New(router infra.Router, db infra.Db, tickenConfig *utils.TickenConfig) *TickenEventApp {
	tickenEventApp := new(TickenEventApp)

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, _ := services.NewProvider(db, tickenConfig)

	tickenEventApp.router = router
	tickenEventApp.serviceProvider = serviceProvider

	var controllers = []api.Controller{
		eventController.NewEventController(serviceProvider),
		organizationController.NewOrganizationController(serviceProvider),
	}

	var appMiddlewares = []api.Middleware{
		middlewares.GetUserMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		router.Use(middleware)
	}

	for _, controller := range controllers {
		controller.Setup(router)
	}

	return tickenEventApp
}

func (tickenEventApp *TickenEventApp) Start() {
	err := tickenEventApp.router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func (tickenEventApp *TickenEventApp) Populate() {
	eventManager := tickenEventApp.serviceProvider.GetEventManager()
	_, err := eventManager.AddEvent("test-event-id", "organizer", "ticken-test-channel")
	if err != nil {
		return // HANDLER DUPLICATES
	}
}
