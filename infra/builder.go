package infra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-event-service/config"
	"ticken-event-service/env"
	"ticken-event-service/infra/bus"
	"ticken-event-service/infra/db"
	"ticken-event-service/infra/hsm"
	"ticken-event-service/log"
	"ticken-event-service/utils"
)

type Builder struct {
	tickenConfig *config.Config
}

var pc peerconnector.PeerConnector = nil

func NewBuilder(tickenConfig *config.Config) (*Builder, error) {
	if tickenConfig == nil {
		return nil, fmt.Errorf("configuration is mandatory")
	}

	builder := new(Builder)
	builder.tickenConfig = tickenConfig

	return builder, nil
}

func (builder *Builder) BuildDb(connString string) Db {
	var tickenDb Db = nil

	switch builder.tickenConfig.Database.Driver {
	case config.MongoDriver:
		tickenDb = db.NewMongoDb()
	default:
		panic(fmt.Errorf("database driver %s not implemented", builder.tickenConfig.Database.Driver))
	}

	err := tickenDb.Connect(connString)
	if err != nil {
		panic(err)
	}

	return tickenDb
}

func (builder *Builder) BuildHSM(encryptingKey string) HSM {
	rootPath, err := utils.GetServiceRootPath()
	if err != nil {
		panic(err)
	}

	localFileSystemHSM, err := hsm.NewLocalFileSystemHSM(encryptingKey, rootPath)
	if err != nil {
		panic(err)
	}
	return localFileSystemHSM
}

func (builder *Builder) BuildEngine() *gin.Engine {
	return gin.Default()
}

func (builder *Builder) BuildPvtbcCaller() *pvtbc.Caller {
	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		panic(err)
	}
	return caller
}

func (builder *Builder) BuildPvtbcListener() *pvtbc.Listener {
	listener, err := pvtbc.NewListener(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		panic(err)
	}
	return listener
}

func (builder *Builder) BuildBusSubscriber(connString string) BusSubscriber {
	var tickenBus BusSubscriber = nil

	driverToUse := builder.tickenConfig.Bus.Driver
	if env.TickenEnv.IsDev() {
		driverToUse = config.DevBusDriver
	}

	switch driverToUse {
	case config.DevBusDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.DevBusDriver)
		tickenBus = bus.NewTickenDevBusSubscriber()
	case config.RabbitMQDriver:
		log.TickenLogger.Info().Msg("using bus subscriber: " + config.RabbitMQDriver)
		tickenBus = bus.NewRabbitMQSubscriber()
	default:
		err := fmt.Errorf("bus driver %s not implemented", builder.tickenConfig.Bus.Driver)
		log.TickenLogger.Panic().Err(err)
	}

	err := tickenBus.Connect(connString, builder.tickenConfig.Bus.Exchange)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("bus subscriber connection established")

	return tickenBus
}

func (builder *Builder) BuildBusPublisher(connString string) BusPublisher {
	var tickenBus BusPublisher = nil

	driverToUse := builder.tickenConfig.Bus.Driver
	if env.TickenEnv.IsDev() {
		driverToUse = config.DevBusDriver
	}

	switch driverToUse {
	case config.DevBusDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.DevBusDriver)
		tickenBus = bus.NewTickenDevBusPublisher()
	case config.RabbitMQDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.RabbitMQDriver)
		tickenBus = bus.NewRabbitMQPublisher()
	default:
		err := fmt.Errorf("bus driver %s not implemented", builder.tickenConfig.Bus.Driver)
		log.TickenLogger.Panic().Err(err)
	}

	err := tickenBus.Connect(connString, builder.tickenConfig.Bus.Exchange)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("bus publisher connection established")

	return tickenBus
}

func buildPeerConnector(config config.PvtbcConfig) peerconnector.PeerConnector {
	if pc != nil {
		return pc
	}

	var pc peerconnector.PeerConnector
	if env.TickenEnv.IsDev() {
		pc = peerconnector.NewDev(config.MspID, "admin")
	} else {
		pc = peerconnector.New(config.MspID, config.CertificatePath, config.PrivateKeyPath)
	}

	err := pc.Connect(config.PeerEndpoint, config.GatewayPeer, config.TLSCertificatePath)
	if err != nil {
		panic(err)
	}

	return pc
}
