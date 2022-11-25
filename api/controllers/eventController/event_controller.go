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
	router.GET("/organizations/:organizationID/events", controller.GetOrganizationEvents)
	router.GET("/organizations/:organizationID/events/:eventID", controller.GetEvent)
	router.POST("/organizations/:organizationID/events", controller.CreateEvent)
}
