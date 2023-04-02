package eventController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/controllers/baseController"
	"ticken-event-service/services"
)

type EventController struct {
	*baseController.BaseController
}

func New(serviceProvider services.IProvider) *EventController {
	return &EventController{BaseController: baseController.New(serviceProvider)}
}

func (controller *EventController) Setup(router gin.IRouter) {
	group := router.Group("/organizations")
	group.POST("/:organizationID/events", controller.CreateEvent)
	group.GET("/:organizationID/events/:eventID", controller.GetEvent)
	group.GET("/:organizationID/events", controller.GetOrganizationEvents)
	group.PATCH("/:organizationID/events/:eventID/on_sale", controller.SetEventOnSale)
	group.PUT("/:organizationID/events/:eventID/sections", controller.AddSection)
}
