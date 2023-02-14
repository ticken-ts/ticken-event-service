package healthController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

const HealthMessage = "Everything is fine"

type HealthController struct {
	serviceProvider services.IProvider
	middlewares     []api.Middleware
}

func New(serviceProvider services.IProvider, middlewares ...api.Middleware) *HealthController {
	controller := new(HealthController)
	controller.serviceProvider = serviceProvider
	controller.middlewares = middlewares
	return controller
}

func (controller *HealthController) Setup(router gin.IRouter) {
	group := router.Group("/metrics")
	for _, middleware := range controller.middlewares {
		middleware.Setup(group)
	}
	group.GET("/healthz", controller.Healthz)
}

func (controller *HealthController) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, utils.HttpResponse{Message: HealthMessage})
}
