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
)

type Builder struct {
	tickenConfig *config.Config
}

var pc *peerconnector.PeerConnector = nil

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

func (builder *Builder) BuildBusPublisher(connString string) BusPublisher {
	var busPublisher BusPublisher = nil

	switch builder.tickenConfig.Bus.Driver {
	case config.RabbitMQDriver:
		busPublisher = bus.NewRabbitMQPublisher()
	default:
		panic(fmt.Errorf("bus driver %s not implemented", builder.tickenConfig.Bus.Driver))
	}

	// if we are on dev, we are running each service separately.
	// So there is no need to use a real. For this case, we are going
	// to use a dev bus that mock all calls
	if env.TickenEnv.IsDev() {
		busPublisher = bus.NewTickenDevBusPublisher()
	}

	err := busPublisher.Connect(connString, builder.tickenConfig.Bus.Exchange)
	if err != nil {
		panic(err)
	}

	return busPublisher
}

func buildPeerConnector(config config.PvtbcConfig) *peerconnector.PeerConnector {
	if pc != nil {
		return pc
	}

	pc := peerconnector.New(config.MspID, config.CertificatePath, config.PrivateKeyPath)

	err := pc.Connect(config.PeerEndpoint, config.GatewayPeer, config.TLSCertificatePath)
	if err != nil {
		panic(err)
	}

	return pc
}
