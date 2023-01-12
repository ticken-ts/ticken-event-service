package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"ticken-event-service/api"
	"ticken-event-service/api/controllers/eventController"
	"ticken-event-service/api/controllers/healthController"
	"ticken-event-service/api/controllers/sectionController"
	"ticken-event-service/api/middlewares"
	"ticken-event-service/api/security"
	"ticken-event-service/app/fakes"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/listeners"
	"ticken-event-service/repos"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

type TickenEventApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider

	// populators are intended to populate
	// useful data. It can be testdata or
	// data that should be present on the db
	// before the service is available
	populators []Populator
}

func New(builder infra.IBuilder, tickenConfig *config.Config) *TickenEventApp {
	tickenEventApp := new(TickenEventApp)

	db := builder.BuildDb(env.TickenEnv.DbConnString)
	hsm := builder.BuildHSM(env.TickenEnv.HSMEncryptionKey)
	engine := builder.BuildEngine()
	pvtbcListener := builder.BuildPvtbcListener()
	busPublisher := builder.BuildBusPublisher(env.TickenEnv.BusConnString)

	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	serviceProvider, err := services.NewProvider(repoProvider, busPublisher, hsm)
	if err != nil {
		panic(err)
	}

	tickenEventApp.engine = engine
	tickenEventApp.config = tickenConfig
	tickenEventApp.repoProvider = repoProvider
	tickenEventApp.serviceProvider = serviceProvider

	var appListeners = []listeners.Listener{
		listeners.NewEventListener(serviceProvider, pvtbcListener, "ticken-channel"),
	}

	for _, listener := range appListeners {
		listener.Listen()
	}

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, &tickenConfig.Server, &tickenConfig.Dev),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	var controllers = []api.Controller{
		eventController.New(serviceProvider),
		healthController.New(serviceProvider),
		sectionController.New(serviceProvider),
	}

	for _, controller := range controllers {
		controller.Setup(engine)
	}

	tickenEventApp.populators = []Populator{
		fakes.NewFakeUsersPopulator(repoProvider.GetOrganizerRepository(), tickenConfig.Dev.User),
		fakes.NewFakeOrgsPopulator(hsm, tickenConfig.Dev.User, repoProvider.GetOrganizerRepository(), repoProvider.GetOrganizationRepository()),
	}

	return tickenEventApp
}

func (tickenEventApp *TickenEventApp) Start() {
	url := tickenEventApp.config.Server.GetServerURL()
	err := tickenEventApp.engine.Run(url)
	if err != nil {
		panic(err)
	}
}

func (tickenEventApp *TickenEventApp) Populate() {
	for _, populator := range tickenEventApp.populators {
		err := populator.Populate()
		if err != nil {
			panic(err)
		}
	}

}

func (tickenEventApp *TickenEventApp) EmitFakeJWT() {
	rsaPrivKey, err := utils.LoadRSA(tickenEventApp.config.Dev.JWTPrivateKey, tickenEventApp.config.Dev.JWTPublicKey)
	if err != nil {
		panic(err)
	}

	fakeJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, &security.Claims{
		Subject:           tickenEventApp.config.Dev.User.UserID,
		Email:             tickenEventApp.config.Dev.User.Email,
		PreferredUsername: tickenEventApp.config.Dev.User.Username,
	})

	signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

	if err != nil {
		panic(fmt.Errorf("error generation fake JWT: %s", err.Error()))
	}

	fmt.Printf("DEV JWT: %s \n", signedJWT)
}
