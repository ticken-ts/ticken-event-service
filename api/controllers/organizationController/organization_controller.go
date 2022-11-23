package organizationController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type OrganizationController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *OrganizationController {
	controller := new(OrganizationController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *OrganizationController) Setup(router gin.IRouter) {
	router.POST("/organizations", controller.RegisterOrganization)
}
