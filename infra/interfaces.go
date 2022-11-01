package infra

import (
	"context"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-event-service/infra/bus"
)

type Db interface {
	Connect(connString string) error
	IsConnected() bool

	// GetClient is going to return the raw client.
	// The caller should cast the returned value
	// into the correct client depending on the
	// driver
	GetClient() interface{}
}

type BusSubscriber interface {
	Connect(connString string, exchangeName string) error
	IsConnected() bool
	Listen(handler func([]byte)) error
}

type BusPublisher interface {
	Connect(connString string, exchangeName string) error
	IsConnected() bool
	Publish(ctx context.Context, msg bus.Message) error
}

type IBuilder interface {
	BuildDb(connString string) Db
	BuildEngine() *gin.Engine
	BuildPvtbcCaller() *pvtbc.Caller
	BuildPvtbcListener() *pvtbc.Listener
	BuildBusPublisher(connString string) BusPublisher
}
