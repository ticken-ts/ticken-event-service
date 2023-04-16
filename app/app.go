package app

import (
	"fmt"
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
	"ticken-event-service/log"
	"ticken-event-service/repos"
	"ticken-event-service/security/jwt"
	"ticken-event-service/services"
	"ticken-event-service/utils"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
)

type TickenEventApp struct {
	engine          *gin.Engine
	jwtVerifier     jwt.Verifier
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
	log.TickenLogger.Info().Msg(color.BlueString("initializing " + tickenConfig.Server.ClientID))

	tickenEventApp := new(TickenEventApp)

	/******************************** infra builds ********************************/
	engine := infraBuilder.BuildEngine()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	fileUploader := infraBuilder.BuildFileUploader()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	hsm := infraBuilder.BuildHSM(env.TickenEnv.HSMEncryptionKey)
	pubbcAdmin := infraBuilder.BuildPubbcAdmin(env.TickenEnv.TickenWalletKey)
	busPublisher := infraBuilder.BuildBusPublisher(env.TickenEnv.BusConnString)
	authIssuer := infraBuilder.BuildAuthIssuer(env.TickenEnv.ServiceClientSecret)
	/**************************++***************************************************/

	/********************************** providers **********************************/
	repoProvider, err := repos.NewProvider(
		db,
		&tickenConfig.Database,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}

	serviceProvider, err := services.NewProvider(
		repoProvider,
		busPublisher,
		hsm,
		infraBuilder,
		pubbcAdmin,
		fileUploader,
		authIssuer,
		tickenConfig,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}
	/**************************++***************************************************/

	tickenEventApp.engine = engine
	tickenEventApp.config = tickenConfig
	tickenEventApp.jwtVerifier = jwtVerifier
	tickenEventApp.repoProvider = repoProvider
	tickenEventApp.serviceProvider = serviceProvider

	tickenEventApp.loadMiddlewares(engine)
	tickenEventApp.loadControllers(engine)

	/********************************* populators **********************************/
	tickenEventApp.populators = []Populator{
		fakes.NewFakeLoader(repoProvider, serviceProvider, tickenConfig, hsm),
	}
	/**************************++***************************************************/

	return tickenEventApp
}

func (tickenEventApp *TickenEventApp) Start() {
	url := tickenEventApp.config.Server.GetServerURL()
	err := tickenEventApp.engine.Run(url)
	if err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to start server: %s", err.Error()))
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
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to load dev RSA: %s", err.Error()))
	}

	organizers := tickenEventApp.repoProvider.GetOrganizerRepository().FindAll()

	for _, organizer := range organizers {
		fakeJWT := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &jwt.Claims{
			Subject:           organizer.OrganizerID.String(),
			Email:             organizer.Email,
			PreferredUsername: organizer.Username,
		})

		signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

		if err != nil {
			log.TickenLogger.Panic().Msg(fmt.Sprintf(
				"error generation fake JWT for user %s: %s", organizer.Username, err.Error()),
			)
		}

		log.TickenLogger.Info().Msg(fmt.Sprintf("dev JWT user: %s -> %s",
			color.GreenString(organizer.Username),
			color.YellowString(signedJWT)))

		if _, err := tickenEventApp.jwtVerifier.Verify(signedJWT); err != nil {
			log.TickenLogger.Warn().Msg(fmt.Sprintf("failed to verify JWT: %s", err.Error()))
		} else {
			log.TickenLogger.Info().Msg(fmt.Sprintf("jwt verified successfully"))
		}

		fmt.Println() // add a new line
	}

}

func (tickenEventApp *TickenEventApp) loadControllers(apiRouter gin.IRouter) {
	apiRouterGroup := apiRouter.Group(tickenEventApp.config.Server.APIPrefix)

	var appControllers = []api.Controller{
		eventController.New(tickenEventApp.serviceProvider),
		healthController.New(tickenEventApp.serviceProvider),
		publicController.New(tickenEventApp.serviceProvider),
		assetController.New(tickenEventApp.serviceProvider),
		validatorController.New(tickenEventApp.serviceProvider),
	}

	for _, controller := range appControllers {
		controller.Setup(apiRouterGroup)
	}
}

func (tickenEventApp *TickenEventApp) loadMiddlewares(apiRouter gin.IRouter) {
	var appMiddlewares = []api.Middleware{
		middlewares.NewCorsMiddleware(),
		middlewares.NewErrorMiddleware(),
		middlewares.NewLoggerMiddleware(),
		middlewares.NewAuthMiddleware(tickenEventApp.jwtVerifier, tickenEventApp.config.Server.APIPrefix),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}
}
