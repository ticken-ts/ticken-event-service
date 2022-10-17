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
	"ticken-event-service/repos"
	"ticken-event-service/services"
)

type TickenEventApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
}

func New(builder infra.IBuilder, tickenConfig *config.Config) *TickenEventApp {
	tickenEventApp := new(TickenEventApp)

	db := builder.BuildDb(env.TickenEnv.ConnString)
	engine := builder.BuildEngine()
	pvtbcListener := builder.BuildPvtbcListener()

	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	serviceProvider, err := services.NewProvider(repoProvider)
	if err != nil {
		panic(err)
	}

	tickenEventApp.engine = engine
	tickenEventApp.repoProvider = repoProvider
	tickenEventApp.serviceProvider = serviceProvider

	var appListeners = []listeners.Listener{
		listeners.NewEventListener(serviceProvider, pvtbcListener, "ticken-channel"),
	}

	for _, listener := range appListeners {
		listener.Listen()
	}

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	var controllers = []api.Controller{
		eventController.NewEventController(serviceProvider),
		organizationController.NewOrganizationController(serviceProvider),
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
