package healthController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

const HealthMessage = "Everything is fine"

type HealthController struct {
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *HealthController {
	controller := new(HealthController)
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *HealthController) Setup(router gin.IRouter) {
	router.GET("/healthz", controller.Healthz)
}

func (controller *HealthController) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, utils.HttpResponse{Message: HealthMessage})
}
