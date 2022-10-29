package listeners

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/cclisteners"
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
	listener.pvtbcListener.Listen(context.Background(), listener.generalEventChaincodeCallback)
}

func (listener *EventListener) generalEventChaincodeCallback(eventNotification *cclisteners.CCEventNotification) {
	switch eventNotification.Type {

	case cclisteners.EventCreatedNotification:
		listener.handleEventCreation(eventNotification)

	case cclisteners.SectionAddedNotification:
		listener.handleSectionAddition(eventNotification)

	default:
		log.TickenLogger.Error().Msgf("event notification type %s not implementes", eventNotification.Type)
	}
}

func (listener *EventListener) handleEventCreation(eventNotification *cclisteners.CCEventNotification) {
	log.TickenLogger.Info().Msgf("handling %s notification: %s", cclisteners.TickenEventChaincode, cclisteners.EventCreatedNotification)
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

func (listener *EventListener) handleSectionAddition(eventNotification *cclisteners.CCEventNotification) {
	log.TickenLogger.Info().Msgf("handling %s notification: %s", cclisteners.TickenEventChaincode, cclisteners.SectionAddedNotification)
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
