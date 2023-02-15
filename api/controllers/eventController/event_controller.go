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
	group := router.Group("/organizations")
	group.POST("/:organizationID/events", controller.CreateEvent)
	group.GET("/:organizationID/events/:eventID", controller.GetEvent)
	group.GET("/:organizationID/events", controller.GetOrganizationEvents)
	group.PATCH("/:organizationID/events/:eventID/on_sale", controller.SetEventOnSale)
	group.PUT("/:organizationID/events/:eventID/sections", controller.AddSection)
}
