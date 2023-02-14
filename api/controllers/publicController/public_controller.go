package publicController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api"
	"ticken-event-service/services"
)

type PublicController struct {
	serviceProvider services.IProvider
	middlewares     []api.Middleware
}

func New(serviceProvider services.IProvider, middlewares ...api.Middleware) *PublicController {
	controller := new(PublicController)
	controller.serviceProvider = serviceProvider
	controller.middlewares = middlewares
	return controller
}

func (controller *PublicController) Setup(router gin.IRouter) {
	group := router.Group("/public")
	for _, middleware := range controller.middlewares {
		middleware.Setup(group)
	}
	group.GET("/events", controller.GetAvailableEvents)
}
