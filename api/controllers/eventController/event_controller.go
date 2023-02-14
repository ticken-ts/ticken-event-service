package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/api"
	"ticken-event-service/services"
)

type EventController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
	middlewares     []api.Middleware
}

func New(serviceProvider services.IProvider, middlewares ...api.Middleware) *EventController {
	controller := new(EventController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	controller.middlewares = middlewares
	return controller
}

func (controller *EventController) Setup(router gin.IRouter) {
	group := router.Group("/organizations")
	for _, middleware := range controller.middlewares {
		middleware.Setup(group)
	}
	group.POST("/:organizationID/events", controller.CreateEvent)
	group.GET("/:organizationID/events/:eventID", controller.GetEvent)
	group.GET("/:organizationID/events", controller.GetOrganizationEvents)
	group.PATCH("/:organizationID/events/:eventID/on_sale", controller.SetEventOnSale)
	group.PUT("/:organizationID/events/:eventID/sections", controller.AddSection)
}
