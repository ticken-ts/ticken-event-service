package publicController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/services"
)

type PublicController struct {
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *PublicController {
	controller := new(PublicController)
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *PublicController) Setup(router gin.IRouter) {
	group := router.Group("/public")
	group.GET("/events", controller.GetAvailableEvents)
	group.GET("/events/:eventID", controller.GetPublicEvent)
}
