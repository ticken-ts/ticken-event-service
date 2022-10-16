package app

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api"
	"ticken-event-service/api/controllers/eventController"
	"ticken-event-service/api/controllers/organizationController"
	"ticken-event-service/api/middlewares"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/listeners"
	"ticken-event-service/services"
)

type TickenEventApp struct {
	engine          *gin.Engine
	serviceProvider services.Provider
}

func New(builder infra.IBuilder, tickenConfig *config.Config) *TickenEventApp {
	tickenEventApp := new(TickenEventApp)

	db := builder.BuildDb(env.TickenEnv.ConnString)
	engine := builder.BuildEngine()
	pvtbcListener := builder.BuildPvtbcListener()

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, _ := services.NewProvider(db, tickenConfig)

	tickenEventApp.engine = engine
	tickenEventApp.serviceProvider = serviceProvider

	var appListeners = []listeners.Listener{
		listeners.NewEventListener(serviceProvider, pvtbcListener, "ticken-channel"),
	}

	var controllers = []api.Controller{
		eventController.NewEventController(serviceProvider),
		organizationController.NewOrganizationController(serviceProvider),
	}

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	for _, listener := range appListeners {
		listener.Listen()
	}

	for _, controller := range controllers {
		controller.Setup(engine)
	}

	return tickenEventApp
}

func (tickenEventApp *TickenEventApp) Start() {
	err := tickenEventApp.engine.Run("localhost:8080")
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
