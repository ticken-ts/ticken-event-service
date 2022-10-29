package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type EventController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *EventController {
	controller := new(EventController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *EventController) Setup(router gin.IRouter) {
	router.GET("/events/:eventId", controller.GetEvent)
	router.GET("/events", controller.GetUserEvents)
	router.POST("/events", controller.CreateEvent)
}
