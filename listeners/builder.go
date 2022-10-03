package listeners

import (
	"context"
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

type Builder struct {
	tickenConfig  *utils.TickenConfig
	pvtbcListener *pvtbc.Listener
}

func NewBuilder(tickenConfig *utils.TickenConfig) (*Builder, error) {
	if tickenConfig == nil {
		return nil, fmt.Errorf("configuration is mandatory")
	}

	builder := new(Builder)
	builder.tickenConfig = tickenConfig

	listener, err := pvtbc.NewListener(
		tickenConfig.Config.Pvtbc.MspID,
		tickenConfig.Config.Pvtbc.CertificatePath,
		tickenConfig.Config.Pvtbc.PrivateKeyPath,
		tickenConfig.Config.Pvtbc.PeerEndpoint,
		tickenConfig.Config.Pvtbc.GatewayPeer,
		tickenConfig.Config.Pvtbc.TLSCertificatePath,
	)
	if err != nil {
		panic(err)
	}

	// TODO: add channel to config
	err = listener.SetChannel(context.Background(), "ticken-channel")

	if err != nil {
		panic(fmt.Errorf("Error setting listener channel: %s", err))
	}

	builder.pvtbcListener = listener

	return builder, nil
}

func (builder *Builder) BuildEventListener(serviceProvider services.Provider) *EventListener {
	a := builder.pvtbcListener
	if a != nil {

	}
	newEventListener := NewEventListener(serviceProvider, builder.pvtbcListener, "ticken-channel")
	return newEventListener
}
