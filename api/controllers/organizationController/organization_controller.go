package organizationController

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"ticken-event-service/infra"
	"ticken-event-service/services"
)

type OrganizationController struct {
	validator       *validator.Validate
	serviceProvider services.Provider
}

// TODO -> test only until user management is complete
var owner = uuid.New().String()

func NewOrganizationController(serviceProvider services.Provider) *OrganizationController {
	controller := new(OrganizationController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *OrganizationController) Setup(router infra.Router) {
	router.GET("/orgs/:organizationId", controller.GetOrganization)
	router.GET("/orgs", controller.GetUserOrganizations)
	//router.PUT("/events/:eventID/tickets/:ticketID/sign", controller.SignTicket) // <- Es REST LCTM
}
