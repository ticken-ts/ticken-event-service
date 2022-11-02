package organizerController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type OrganizerController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *OrganizerController {
	controller := new(OrganizerController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *OrganizerController) Setup(router gin.IRouter) {
	router.POST("/organizers", controller.RegisterOrganizer)
}
