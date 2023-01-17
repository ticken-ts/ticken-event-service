package sectionController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type SectionController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *SectionController {
	controller := new(SectionController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *SectionController) Setup(router gin.IRouter) {
	router.PUT("/organizations/:organizationID/events/:eventID/sections", controller.AddSection)
}
