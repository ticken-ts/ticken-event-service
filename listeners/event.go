package listeners

import (
	"github.com/go-playground/validator/v10"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-event-service/services"
)

type EventListener struct {
	validator       *validator.Validate
	serviceProvider services.Provider
	pvtbcListener   *pvtbc.Listener
	channel         string
}

func NewEventListener(serviceProvider services.Provider, pvtbcListener *pvtbc.Listener, channel string) *EventListener {
	newEventListener := new(EventListener)
	newEventListener.pvtbcListener = pvtbcListener
	newEventListener.serviceProvider = serviceProvider
	newEventListener.channel = channel
	return newEventListener
}

func (listener *EventListener) Listen() {

	callback := func(event *chain_models.Event) {
		_, err := listener.serviceProvider.GetEventManager().AddEvent(
			event.EventID,
			event.OrganizationID,
			listener.channel,
		)
		if err != nil {
			panic("error adding pvtbc event")
		}
	}

	err := listener.pvtbcListener.ListenNewEvents(callback)

	if err != nil {
		panic(err)
	}
}
