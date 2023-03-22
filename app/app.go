package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
	"ticken-event-service/api"
	"ticken-event-service/api/controllers/assetController"
	"ticken-event-service/api/controllers/eventController"
	"ticken-event-service/api/controllers/healthController"
	"ticken-event-service/api/controllers/publicController"
	"ticken-event-service/api/controllers/validatorController"
	"ticken-event-service/api/middlewares"
	"ticken-event-service/app/fakes"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra"
	"ticken-event-service/repos"
	"ticken-event-service/security/jwt"
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

func New(infraBuilder infra.IBuilder, tickenConfig *config.Config) *TickenEventApp {
	tickenEventApp := new(TickenEventApp)

	engine := infraBuilder.BuildEngine()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	fileUploader := infraBuilder.BuildFileUploader()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	hsm := infraBuilder.BuildHSM(env.TickenEnv.HSMEncryptionKey)
	pubbcAdmin := infraBuilder.BuildPubbcAdmin(env.TickenEnv.TickenWalletKey)
	busPublisher := infraBuilder.BuildBusPublisher(env.TickenEnv.BusConnString)
	authIssuer := infraBuilder.BuildAuthIssuer(env.TickenEnv.ServiceClientSecret)

	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	serviceProvider, err := services.NewProvider(
		repoProvider,
		busPublisher,
		hsm,
		infraBuilder,
		pubbcAdmin,
		fileUploader,
		authIssuer,
		tickenConfig.Services,
	)
	if err != nil {
		panic(err)
	}

	tickenEventApp.engine = engine
	tickenEventApp.config = tickenConfig
	tickenEventApp.repoProvider = repoProvider
	tickenEventApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, jwtVerifier, tickenConfig.Server.APIPrefix),
	}

	var controllers = []api.Controller{
		eventController.New(serviceProvider),
		healthController.New(serviceProvider),
		publicController.New(serviceProvider),
		assetController.New(serviceProvider),
		validatorController.New(serviceProvider),
	}

	apiRouter := engine.Group(tickenConfig.Server.APIPrefix)

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}

	for _, controller := range controllers {
		controller.Setup(apiRouter)
	}

	tickenEventApp.populators = []Populator{
		fakes.NewFakeUsersPopulator(repoProvider, authIssuer, tickenConfig.Dev, tickenConfig.Services),
		fakes.NewFakeOrgsPopulator(repoProvider, authIssuer, tickenConfig.Dev, hsm, tickenConfig.Pvtbc.ClusterStoragePath),
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

	fakeJWT := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &jwt.Claims{
		Subject:           tickenEventApp.config.Dev.User.UserID,
		Email:             tickenEventApp.config.Dev.User.Email,
		PreferredUsername: tickenEventApp.config.Dev.User.Username,
	})

	signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

	if err != nil {
		panic(fmt.Errorf("error generation fake Token: %s", err.Error()))
	}

	fmt.Printf("DEV Token: %s \n", signedJWT)
}
