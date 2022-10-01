package eventController

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"ticken-event-service/infra"
	"ticken-event-service/services"
)

type EventController struct {
	validator       *validator.Validate
	serviceProvider services.Provider
}

// TODO -> test only until user management is complete
var owner = uuid.New().String()

func NewEventController(serviceProvider services.Provider) *EventController {
	controller := new(EventController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *EventController) Setup(router infra.Router) {
	router.GET("/events/:eventID", controller.GetEvent)
	router.GET("/events", controller.GetUserEvents)
	//router.PUT("/events/:eventID/tickets/:ticketID/sign", controller.SignTicket) // <- Es REST LCTM
}
