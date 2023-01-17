package listeners

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"ticken-event-service/log"
	"ticken-event-service/services"
)

type EventListener struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
	pvtbcListener   *pvtbc.Listener
	channel         string
}

func NewEventListener(serviceProvider services.IProvider, pvtbcListener *pvtbc.Listener, channel string) *EventListener {
	newEventListener := new(EventListener)

	newEventListener.channel = channel
	newEventListener.pvtbcListener = pvtbcListener
	newEventListener.serviceProvider = serviceProvider

	err := pvtbcListener.SetChannel(channel)
	if err != nil {
		panic(err)
	}

	return newEventListener
}

func (listener *EventListener) Listen() {
	listener.pvtbcListener.ListenCCEvent(context.Background(), listener.generalEventChaincodeCallback)
}

func (listener *EventListener) generalEventChaincodeCallback(eventNotification *ccclient.CCNotification) {
	switch eventNotification.Type {

	case consts.EventCreatedNotification:
		listener.handleEventCreation(eventNotification)

	case consts.SectionAddedNotification:
		listener.handleSectionAddition(eventNotification)

	default:
		log.TickenLogger.Error().Msgf("event notification type %s not implementes", eventNotification.Type)
	}
}

func (listener *EventListener) handleEventCreation(eventNotification *ccclient.CCNotification) {
	log.TickenLogger.Info().Msgf("handling %s notification: %s", consts.TickenEventChaincode, consts.EventCreatedNotification)
	eventManager := listener.serviceProvider.GetEventManager()

	var onChainEvent chainmodels.Event
	err := json.Unmarshal(eventNotification.Payload, &onChainEvent)
	if err != nil {
		log.TickenLogger.Error().Err(err)
		return
	}

	_, err = eventManager.SyncOnChainEvent(&onChainEvent, listener.channel)
	if err != nil {
		log.TickenLogger.Error().Err(err)
		return
	}
	log.TickenLogger.Info().Msgf("on chain event %s sync", onChainEvent.EventID)
}

func (listener *EventListener) handleSectionAddition(eventNotification *ccclient.CCNotification) {
	log.TickenLogger.Info().Msgf("handling %s notification: %s", consts.TickenEventChaincode, consts.SectionAddedNotification)
	eventManager := listener.serviceProvider.GetEventManager()

	var onChainSection chainmodels.Section
	err := json.Unmarshal(eventNotification.Payload, &onChainSection)
	if err != nil {
		log.TickenLogger.Error().Err(err)
		return
	}

	_, err = eventManager.SyncOnChainSection(&onChainSection)
	if err != nil {
		log.TickenLogger.Error().Err(err)
		return
	}
	log.TickenLogger.Info().Msgf("on chain section %s sync for event %s", onChainSection.Name, "")

}
