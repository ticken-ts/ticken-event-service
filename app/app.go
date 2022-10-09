package app

import (
	"ticken-event-service/api"
	"ticken-event-service/api/controllers/eventController"
	"ticken-event-service/api/controllers/organizationController"
	"ticken-event-service/api/middlewares"
	"ticken-event-service/infra"
	"ticken-event-service/listeners"
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

	listenerBuilder, _ := listeners.NewBuilder(tickenConfig)

	tickenEventApp.router = router
	tickenEventApp.serviceProvider = serviceProvider

	var appListeners = []listeners.Listener{
		listenerBuilder.BuildEventListener(serviceProvider),
	}

	var controllers = []api.Controller{
		eventController.NewEventController(serviceProvider),
		organizationController.NewOrganizationController(serviceProvider),
	}

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(router)
	}

	for _, listener := range appListeners {
		listener.Listen()
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
	orgManager := tickenEventApp.serviceProvider.GetOrgManager()
	_, err := orgManager.AddOrganization("organizer", []string{}, []string{"aishd98y8954j5k4m"})
	_, err = eventManager.AddEvent("test-event-id", "organizer", "ticken-test-channel")
	if err != nil {
		return // HANDLER DUPLICATES
	}
}
