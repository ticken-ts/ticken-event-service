package validatorController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type ValidatorController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *ValidatorController {
	controller := new(ValidatorController)
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *ValidatorController) Setup(router gin.IRouter) {
	group := router.Group("/organizations/:organizationID")
	group.POST("/validators", controller.RegisterValidator)
}
